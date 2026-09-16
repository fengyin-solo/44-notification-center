// Package app 负责依赖装配。
package app

import (
	"net/http"

	"notification/internal/config"
	"notification/internal/handler"
	"notification/internal/service"
	"notification/internal/store"
	"notification/pkg/logger"
)

// App 应用实例，聚合配置、存储、服务与路由。
type App struct {
	server *handler.Server
	store  store.Store
	svc    *service.Service
}

// New 装配全部依赖并返回应用。
func New(cfg *config.Config, log *logger.Logger) (*App, error) {
	st := store.NewMemoryStore()
	svc := service.New(st, log, cfg)
	server := handler.NewServer(svc, log, cfg)
	log.Infof("应用装配完成，配置：%s", cfg.String())
	return &App{server: server, store: st, svc: svc}, nil
}

// Routes 返回根路由处理器。
func (a *App) Routes() http.Handler { return a.server.Routes() }
