package main

import (
	"log"

	"github.com/Zrossiz/finance-backend/internal/app"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Default().Fatal(err)
		}
	}()

	app, err := app.New()
	if err != nil {
		logrus.Fatal(err)
	}

	if err := app.Start(); err != nil {
		logrus.Fatal(err)
	}

	if err := app.ErrGroup.Wait(); err != nil {
		logrus.Fatal(err)
	}

	if err := app.Stop(); err != nil {
		logrus.Fatal(err)
	}
}
