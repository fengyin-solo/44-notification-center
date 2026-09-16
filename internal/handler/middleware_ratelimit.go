package handler

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"notification/pkg/httpx"
)

// tokenBucket 简单令牌桶限流器。
type tokenBucket struct {
	mu       sync.Mutex
	rate     float64 // 每秒补充令牌数
	capacity float64 // 桶容量
	tokens   float64 // 当前令牌数
	last     time.Time
}

// newTokenBucket 创建令牌桶。
func newTokenBucket(rate, capacity float64) *tokenBucket {
	return &tokenBucket{
		rate:     rate,
		capacity: capacity,
		tokens:   capacity,
		last:     time.Now(),
	}
}

// allow 尝试取走一个令牌，成功返回 true。
func (b *tokenBucket) allow() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.tokens += elapsed * b.rate
	if b.tokens > b.capacity {
		b.tokens = b.capacity
	}
	b.last = now
	if b.tokens >= 1 {
		b.tokens--
		return true
	}
	return false
}

// rateLimitMiddleware 对 /api/ 路径做令牌桶限流。
func (s *Server) rateLimitMiddleware(next http.Handler) http.Handler {
	rps := 100.0
	burst := 200.0
	if s.cfg != nil {
		if s.cfg.RateLimitRPS > 0 {
			rps = float64(s.cfg.RateLimitRPS)
		}
		if s.cfg.RateBurst > 0 {
			burst = float64(s.cfg.RateBurst)
		}
	}
	bucket := newTokenBucket(rps, burst)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, apiPrefix) && !bucket.allow() {
			httpx.Error(w, http.StatusTooManyRequests, 429, "请求过于频繁，请稍后再试")
			return
		}
		next.ServeHTTP(w, r)
	})
}
