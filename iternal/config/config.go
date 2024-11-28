package config

import (
	"strconv"

	"github.com/caarlos0/env"
	"github.com/labstack/echo/v4/middleware"
	"github.com/pkg/errors"
)

type Config struct {
	*app
	*storage
	*cors
}

func New() (Config, error) {
	app, err := newApp()
	if err != nil {
		return Config{}, errors.Wrap(err, "Config.New")
	}

	storage, err := newStorage()
	if err != nil {
		return Config{}, errors.Wrap(err, "Config.New")
	}

	cors, err := newCORS()
	if err != nil {
		return Config{}, errors.Wrap(err, "Config.New")
	}

	return Config{app, storage, cors}, nil
}

func newStorage() (*storage, error) {
	conf := storage{}

	if err := env.Parse(&conf); err != nil {
		return &storage{},
			errors.Wrap(err, "newStorage.Parse")
	}

	return &conf, nil
}

type storage struct {
	StoragePort     string `env:"STORAGE_PORT" envDefault:"5432"`
	StorageAddr     string `env:"STORAGE_ADDR" envDefault:"localhost"`
	StorageBaseName string `env:"STORAGE_NAME_BASE" envDefault:"test"`
	StorageLogin    string `env:"STORAGE_LOGIN" envDefault:"postgres"`
	StoragePass     string `env:"STORAGE_PASS" envDefault:"pass"`
}

func (s *storage) GetStoragePort() string      { return s.StoragePort }
func (s *storage) GetStorageAddr() string      { return s.StorageAddr }
func (s *storage) GetStorageBaseName() string  { return s.StorageBaseName }
func (s *storage) GetStorageUserLogin() string { return s.StorageLogin }
func (s *storage) GetStorageUserPass() string  { return s.StoragePass }

type app struct {
	AppHTTPPort string `env:"APP_HTTP_PORT" envDefault:"4000"`
	NameHTTP    string `env:"APP_HTTP_NAME" envDefault:"172.25.162.224"`
	LogLevel    int    `env:"LOG_LEVEL" envDefault:"0"`
}

func (a *app) GetAppHTTPPort() string  { return a.AppHTTPPort }
func (a *app) GetNameHTTPPort() string { return a.NameHTTP }
func (a *app) GetLogLevel() int        { return a.LogLevel }

func newApp() (*app, error) {
	conf := app{}

	if err := env.Parse(&conf); err != nil {
		return &app{},
			errors.Wrap(err, "newApp.Parse")
	}

	return &conf, nil
}

func newCORS() (*cors, error) {
	conf := cors{}

	if err := env.Parse(&conf); err != nil {
		return &cors{},
			errors.Wrap(err, "newCORS.Parse")
	}

	return &conf, nil
}

type cors struct {
	CORSAllowOrigins  []string `env:"CORS_ALLOW_ORIGINS"`
	CORSAllowHeaders  []string `env:"CORS_AIIOW_HEADERS"`
	CORSAllowMethods  []string `env:"CORS_AIIOW_METHODS"`
	CORSExposeHeaders []string `env:"CORS_EXPOSE_HEADERS"`
	CORSMaxAge        string   `env:"CORS_MAX_AGE" envDefault:"0"`
}

func (c *cors) GetCORSAllowOrigins() []string {
	if len(c.CORSAllowOrigins) == 0 {
		return middleware.DefaultCORSConfig.AllowOrigins
	}
	return c.CORSAllowOrigins
}

func (c *cors) GetCORSAllowHeaders() []string {
	if len(c.CORSAllowHeaders) == 0 {
		return middleware.DefaultCORSConfig.AllowHeaders
	}
	return c.CORSAllowHeaders
}

func (c *cors) GetCORSAllowMethods() []string {
	if len(c.CORSAllowMethods) == 0 {
		return middleware.DefaultCORSConfig.AllowMethods
	}
	return c.CORSAllowMethods
}

func (c *cors) GetCORSExposeHeaders() []string {
	if len(c.CORSExposeHeaders) == 0 {
		return middleware.DefaultCORSConfig.ExposeHeaders
	}
	return c.CORSExposeHeaders
}

func (c *cors) GetCORSMaxAge() int {
	integer, err := strconv.Atoi(c.CORSMaxAge)
	if err != nil {
		return 0
	}
	return integer
}
