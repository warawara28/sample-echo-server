package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/warawara28/sample-books/app"
)

const (
	timeoutSecond          = time.Second * 1
	gracefulShutdownSecond = time.Second * 10
)

var version string

func main() {
	logger := app.NewZerolog()
	logger.Info("Version:%s", version)

	env := os.Getenv("ENV")
	if env == "" {
		logger.Fatal("Please set environment variable 'ENV'")
	}

	if err := app.LoadConfig(env); err != nil {
		logger.Fatal("Failed to load file '%s' error:%w", env, err)
	}
	logger.Info("Load config file '%s'", env)

	conf := app.Config()

	dsn := conf.GetString("db.dsn")
	db, err := app.NewDatabase(dsn, logger)
	if err != nil {
		logger.Fatal("Failed to connect database '%s' error:%w", dsn, err)
	}
	if conf.GetBool("db.migrate") {
		if err := db.AutoMigrate(
			app.Book{},
		); err != nil {
			logger.Fatal("Failed to migrate database schema error:%w", err)
		}
		logger.Info("Complete database migration")
	}

	srv := app.NewService(db)

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Use(
		app.NewZerologMiddleware(logger),
		middleware.Recover(),
		middleware.Gzip(),
		middleware.TimeoutWithConfig(middleware.TimeoutConfig{Timeout: timeoutSecond}),
	)

	e.GET("healthz", srv.CheckHealth())
	v1 := e.Group("/v1")
	{
		v1.GET("/books", srv.ListBooks())
	}

	address := conf.GetString("app.address")
	logger.Info("http server started on %s", address)
	go func() {
		if err := e.Start(address); err != nil && err != http.ErrServerClosed {
			e.Logger.Fatal("shutting down the server")
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(
		quit,
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
	)

	sig := <-quit
	ctx, cancel := context.WithTimeout(context.Background(), gracefulShutdownSecond)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
	logger.Info("Received %v signal. shutdown server", sig)
}
