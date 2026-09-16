// Package handler 实现 HTTP 处理器层。
package handler

import (
	"errors"
	"net/http"

	"notification/internal/config"
	"notification/internal/model"
	"notification/internal/service"
	"notification/internal/store"
	"notification/pkg/httpx"
	"notification/pkg/logger"
)

// Server 聚合服务、日志与配置，注册全部路由。
type Server struct {
	svc *service.Service
	log *logger.Logger
	cfg *config.Config
}

// NewServer 创建处理器实例。
func NewServer(svc *service.Service, log *logger.Logger, cfg *config.Config) *Server {
	return &Server{svc: svc, log: log, cfg: cfg}
}

// Routes 组装根路由：注册全部 API 与静态页面，并包裹中间件。
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	s.registerChannelRoutes(mux)
	s.registerTemplateRoutes(mux)
	s.registerTopicRoutes(mux)
	s.registerRecipientRoutes(mux)
	s.registerSubscriptionRoutes(mux)
	s.registerMessageRoutes(mux)
	s.registerSendRecordRoutes(mux)
	s.registerRetryPolicyRoutes(mux)
	s.registerScheduleRoutes(mux)
	s.registerStatsRoutes(mux)
	s.registerExportRoutes(mux)

	mux.Handle("GET /", http.FileServer(http.Dir("web")))

	var h http.Handler = mux
	h = s.rateLimitMiddleware(h)
	h = s.authMiddleware(h)
	h = s.loggingMiddleware(h)
	h = s.recoveryMiddleware(h)
	return h
}

// maxPageSize 返回分页上限。
func (s *Server) maxPageSize() int {
	if s.cfg != nil && s.cfg.MaxPageSize > 0 {
		return s.cfg.MaxPageSize
	}
	return 100
}

// writeServiceError 将业务错误映射为 HTTP 响应。
func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case model.IsValidationError(err):
		httpx.BadRequest(w, err.Error())
	case errors.Is(err, store.ErrNotFound):
		httpx.NotFound(w, err.Error())
	case errors.Is(err, store.ErrConflict):
		httpx.Conflict(w, err.Error())
	default:
		httpx.InternalError(w, err.Error())
	}
}
