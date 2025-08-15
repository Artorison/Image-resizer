package app

import (
	"log/slog"

	"github.com/Artorison/Image-resizer/internal/config"
	mw "github.com/Artorison/Image-resizer/internal/middleware"
	"github.com/Artorison/Image-resizer/pkg/logger"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Handlers interface {
	GetImage(c echo.Context) error
}

type App struct {
	Router   *echo.Echo
	Cfg      *config.Config
	Log      logger.Logger
	Handlers Handlers
}

func New(
	router *echo.Echo,
	cfg *config.Config,
	log logger.Logger,
	handlers Handlers,
) *App {
	return &App{
		Router:   router,
		Cfg:      cfg,
		Log:      log,
		Handlers: handlers,
	}
}

func (a *App) Start() {
	a.Log.Info("starting server", slog.String("port", a.Cfg.App.Port))
	if err := a.Router.Start(":" + a.Cfg.App.Port); err != nil {
		a.Log.Error("Server stopped", logger.Err(err))
	}
}

func (a *App) RegisterRoutes() {
	a.Router.Use(mw.Recover(), mw.LoggingMV())
	a.Router.GET("/fill/:width/:height/*", a.Handlers.GetImage)
}

func (a *App) SetupFront() {
	a.Router.HideBanner = true
	a.Router.Pre(middleware.RemoveTrailingSlash())

	a.Router.Static("/", "web")
	a.Router.File("/", "web/index.html")
}
