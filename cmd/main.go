package main

import (
	"context"
	"flag"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/Artorison/Image-resizer/internal/app"
	"github.com/Artorison/Image-resizer/internal/config"
	"github.com/Artorison/Image-resizer/internal/handlers"
	"github.com/Artorison/Image-resizer/internal/services"
	"github.com/Artorison/Image-resizer/pkg/downloader"
	"github.com/Artorison/Image-resizer/pkg/imageprocessor"
	"github.com/Artorison/Image-resizer/pkg/logger"
	lrucache "github.com/Artorison/Image-resizer/pkg/lru_cache"
	"github.com/labstack/echo/v4"
)

func main() {
	configFilepath := flag.String("config", "", "Path to configuration file")
	flag.Parse()
	cfg := config.MustLoad(*configFilepath)

	logg := logger.NewLogger(cfg.Logger.Level, cfg.Logger.Type)
	logg.Debug("config", *cfg)

	cache := lrucache.NewCache(cfg.App.CacheSize)

	client := &http.Client{
		Timeout: cfg.App.Timeout,
	}
	downloader := downloader.New(client, logg)
	imageProc := imageprocessor.New(cfg.App.Quality)
	service := services.New(cache, downloader, imageProc)
	handler := handlers.New(service, logg)

	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	router := echo.New()
	appl := app.New(router, cfg, logg, handler)
	appl.RegisterRoutes()
	go appl.Start()

	<-ctx.Done()
	logg.Info("server graseful stopped...")
}
