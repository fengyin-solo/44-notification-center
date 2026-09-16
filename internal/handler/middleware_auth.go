package handler

import (
	"crypto/subtle"
	"net/http"
	"strings"

	"notification/pkg/httpx"
)

// apiPrefix 鉴权与限流作用的路径前缀。
const apiPrefix = "/api/"

// authMiddleware 校验 X-API-Key，仅对 /api/ 路径生效。
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.HasPrefix(r.URL.Path, apiPrefix) {
			next.ServeHTTP(w, r)
			return
		}
		key := r.Header.Get("X-API-Key")
		expected := ""
		if s.cfg != nil {
			expected = s.cfg.APIKey
		}
		if expected == "" || subtle.ConstantTimeCompare([]byte(key), []byte(expected)) != 1 {
			httpx.Unauthorized(w, "无效的 API Key")
			return
		}
		next.ServeHTTP(w, r)
	})
}
