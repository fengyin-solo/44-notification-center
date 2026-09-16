// Package service 实现业务逻辑层。
package service

import (
	"notification/internal/config"
	"notification/internal/store"
	"notification/pkg/logger"
)

// Service 聚合存储、日志与配置，提供全部业务方法。
type Service struct {
	store store.Store
	log   *logger.Logger
	cfg   *config.Config
}

// New 创建服务实例。
func New(st store.Store, log *logger.Logger, cfg *config.Config) *Service {
	return &Service{store: st, log: log, cfg: cfg}
}

// paginate 对切片按页截取，start 越界时返回空切片。
func paginate[T any](items []T, page, size int) []T {
	start := (page - 1) * size
	if start >= len(items) {
		return []T{}
	}
	end := start + size
	if end > len(items) {
		end = len(items)
	}
	return items[start:end]
}
