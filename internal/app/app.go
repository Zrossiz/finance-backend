package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Zrossiz/finance-backend/internal/config"
	"github.com/Zrossiz/finance-backend/internal/handler"
	pgrepo "github.com/Zrossiz/finance-backend/internal/repository/pg"
	"github.com/Zrossiz/finance-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
)

type App struct {
	Config *config.Config

	server http.Server

	conn *sql.DB

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

	pgRepo := pgrepo.New(app.conn)

	srv := service.New(service.Postgres{
		User:           pgRepo.User,
		RealEstate:     pgRepo.RealEstate,
		CryptoPosition: pgRepo.CryptoPosition,
		Bond:           pgRepo.Bond,
		Stock:          pgRepo.Stock,
		BankDeposit:    pgRepo.BankDeposit,
	}, cfg)

	r := chi.NewRouter()

	httpHandler := handler.New(handler.Service{
		User:           srv.User,
		RealEstate:     srv.RealEstate,
		CryptoPosition: srv.CryptoPosition,
		Bond:           srv.Bond,
		Stock:          srv.Stock,
		BankDeposit:    srv.BankDeposit,
	}, cfg)
	httpHandler.RegisterRoutes(r)

	app.server = http.Server{
		Addr:    cfg.Server.Addr,
		Handler: r,
	}

	return app, nil
}

func (a *App) Start() error {
	logrus.Info("starting http server...")
	a.ErrGroup.Go(func() error {
		err := a.server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("listen and serve err: %w", err)
		}

		return nil
	})
	logrus.Info("http server is ready")

	return nil
}

func (a *App) Stop() error {
	a.cancelFunc()
	err := a.conn.Close()
	if err != nil {
		return err
	}

	return nil
}
