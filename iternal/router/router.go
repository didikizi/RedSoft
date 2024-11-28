package router

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/didikizi/RedSoft/iternal/service"
	utils "github.com/didikizi/RedSoft/packege"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type Router struct {
	Config
	Service
}

type Config interface {
	GetAppHTTPPort() string
	GetCORSAllowOrigins() []string
	GetCORSAllowHeaders() []string
	GetCORSAllowMethods() []string
	GetCORSExposeHeaders() []string
	GetCORSMaxAge() int
}

type Service interface {
	GetHumanList(context.Context) ([]*service.Human, error)
	GetHumanFromId(context.Context, int) ([]*service.Human, error)
	GetHumanFromSurname(context.Context, string) ([]*service.Human, error)
	DeleteHuman(context.Context, int) (bool, error)
	PutHuman(context.Context, service.PutHuman) error
	PostHuman(context.Context, service.PostHuman, int) error

	PutMailForHuman(context.Context, service.PutMail) error
	GetMailListForHuman(context.Context, int) ([]*service.Mail, error)
	DeleteMailForHuman(context.Context, int) (bool, error)
}

func New(config Config, service Service) *Router {
	return &Router{Config: config, Service: service}
}

func (r *Router) Start() {
	server := echo.New()

	server.Use(echoprometheus.NewMiddleware("test_task"))
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:  r.GetCORSAllowOrigins(),
		AllowHeaders:  r.GetCORSAllowHeaders(),
		AllowMethods:  r.GetCORSAllowMethods(),
		ExposeHeaders: r.GetCORSExposeHeaders(),
		MaxAge:        r.GetCORSMaxAge(),
	}))
	server.Use(middleware.Logger())
	server.Use(middleware.Recover())

	server.GET("/metrics", echoprometheus.NewHandler())

	api := server.Group("/api")
	v1 := api.Group("/v1")

	v1.GET("/humans", r.GetHumanList)
	v1.GET("/human/:cursor", r.GetHuman)
	v1.DELETE("/human/:id", r.DeleteHuman)
	v1.PUT("/human", r.PutHuman)
	v1.POST("/human/:id", r.PostHuman)

	v1.PUT("/human/:id/mail", r.PutMailForHuman)
	v1.GET("/human/:id/mail", r.GetMailListForHuman)
	v1.DELETE("/human/:id_human/mail/:id_mail", r.DeleteMailForHuman)

	start(server, r.GetAppHTTPPort())
	shutdown(server)
}

func (r *Router) Ping(ctx echo.Context) error {
	return ctx.NoContent(http.StatusOK)
}

func start(server *echo.Echo, port string) {
	if err := server.Start((":" + port)); err != nil && errors.Is(err, http.ErrServerClosed) {
		slog.Info(utils.GetCallerInfo(), slog.String("start err:", err.Error()))
		os.Exit(1)
	}
}

func shutdown(server *echo.Echo) {
	quite := make(chan os.Signal, 1)
	signal.Notify(quite, syscall.SIGINT, syscall.SIGTERM)
	<-quite

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	if err := server.Shutdown(ctx); err != nil {
		slog.Info(utils.GetCallerInfo(), slog.String("shutdown err:", err.Error()))
	}
	defer cancel()
}
