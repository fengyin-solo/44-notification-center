// Package config 负责从环境变量加载服务配置。
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config 服务运行配置。
type Config struct {
	Addr         string
	MaxPageSize  int
	APIKey       string
	RateLimitRPS int
	RateBurst    int
}

// Load 从环境变量读取配置，缺失时回退默认值。
func Load() *Config {
	cfg := &Config{
		Addr:         ":" + getenv("PORT", "8080"),
		MaxPageSize:  getenvInt("MAX_PAGE_SIZE", 100),
		APIKey:       getenv("API_KEY", "dev-notify-key"),
		RateLimitRPS: getenvInt("RATE_LIMIT_RPS", 100),
		RateBurst:    getenvInt("RATE_BURST", 200),
	}
	if v := os.Getenv("ADDR"); v != "" {
		cfg.Addr = v
	}
	return cfg
}

func getenv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return strings.TrimSpace(v)
	}
	return def
}

func getenvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil || n <= 0 {
		return def
	}
	return n
}

func (c *Config) String() string {
	return fmt.Sprintf("addr=%s max_page_size=%d rate_limit_rps=%d", c.Addr, c.MaxPageSize, c.RateLimitRPS)
}
