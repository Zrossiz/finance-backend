package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Zrossiz/finance-backend/internal/api"
	"github.com/Zrossiz/finance-backend/internal/config"
	cronscheduler "github.com/Zrossiz/finance-backend/internal/cron_scheduler"
	"github.com/Zrossiz/finance-backend/internal/handler"
	"github.com/Zrossiz/finance-backend/internal/repository/cache"
	pgrepo "github.com/Zrossiz/finance-backend/internal/repository/pg"
	"github.com/Zrossiz/finance-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Config *config.Config

	server http.Server

	conn    *sql.DB
	rdbConn *redis.Client

	CronScheduler *cronscheduler.CronScheduler

	Context    context.Context
	ErrGroup   *errgroup.Group
	cancelFunc context.CancelFunc
}

func New() (*App, error) {
	app := &App{}

	cfg, err := config.New()
	if err != nil {
		return nil, err
	}
	logrus.SetLevel(logrus.Level(cfg.App.Severity))

	app.Config = cfg
	app.Context, app.cancelFunc = signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	app.ErrGroup, app.Context = errgroup.WithContext(app.Context)

	app.conn, err = pgrepo.Connect(fmt.Sprintf(
		"postgresql://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.DB,
	))
	if err != nil {
		return nil, err
	}
	logrus.Info("successfull postgres connect")

	pgRepo := pgrepo.New(app.conn)

	apiSrv := api.NewAPI()

	app.rdbConn, err = cache.Connect(cfg)
	if err != nil {
		return nil, err
	}
	logrus.Info("successfull redis connect")

	rdbCache := cache.New(app.rdbConn)

	srv, err := service.New(service.Postgres{
		User:           pgRepo.User,
		RealEstate:     pgRepo.RealEstate,
		CryptoPosition: pgRepo.CryptoPosition,
		Bond:           pgRepo.Bond,
		Stock:          pgRepo.Stock,
		BankDeposit:    pgRepo.BankDeposit,
	}, service.API{
		CryptoRates: apiSrv.CryptoRates,
	}, service.Cache{
		CryptoRates: &rdbCache.CryptoRates,
	}, cfg)
	if err != nil {
		return nil, err
	}

	r := chi.NewRouter()

	var allowedOrigins []string
	for _, origin := range strings.Split(cfg.Server.CORSAllowedOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			allowedOrigins = append(allowedOrigins, origin)
		}
	}

	if len(allowedOrigins) == 0 {
		if cfg.App.ENV == "prod" {
			return nil, errors.New("CORS allowed origins are empty in production")
		}

		allowedOrigins = []string{"http://localhost:5173"}
	}

	logrus.Info("allowed origins: ", allowedOrigins)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	httpHandler, err := handler.New(handler.Service{
		User:           srv.User,
		RealEstate:     srv.RealEstate,
		CryptoPosition: srv.CryptoPosition,
		Bond:           srv.Bond,
		Stock:          srv.Stock,
		BankDeposit:    srv.BankDeposit,
	}, cfg)
	if err != nil {
		return nil, err
	}

	httpHandler.RegisterRoutes(r, cfg.Server.JWTAccessSecret)

	app.server = http.Server{
		Addr:              cfg.Server.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	err = app.StartCron(srv)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (a *App) Start() error {
	logrus.Info("starting http server...")

	a.ErrGroup.Go(func() error {
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve: %w", err)
		}

		return nil
	})

	a.ErrGroup.Go(func() error {
		<-a.Context.Done()

		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			10*time.Second,
		)
		defer cancel()

		if err := a.server.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutdown http server: %w", err)
		}

		return nil
	})

	logrus.Info("http server is ready")

	return nil
}

func (a *App) StartCron(srv *service.Service) error {
	cr := cron.New()

	a.CronScheduler = cronscheduler.New(cr)
	err := a.CronScheduler.AddTask(cronscheduler.Task{
		Name:     "refresh_crypto_rates",
		Schedule: "* * * * *",
		Handler:  srv.CryptoRates.RefreshCryptoRatesCache,
	})
	if err != nil {
		return err
	}

	cr.Start()

	logrus.Infof("cron started, total tasks: %v", len(a.CronScheduler.Tasks))

	a.CronScheduler.StartOnInit(a.CronScheduler.Tasks)

	return nil
}

func (a *App) Stop() error {
	a.cancelFunc()

	ctx := a.CronScheduler.Stop()
	<-ctx.Done()

	if err := a.rdbConn.Close(); err != nil {
		return fmt.Errorf("close redis: %w", err)
	}

	if err := a.conn.Close(); err != nil {
		return fmt.Errorf("close postgres: %w", err)
	}

	return nil
}
