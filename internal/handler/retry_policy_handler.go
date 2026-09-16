package handler

import (
	"net/http"

	"notification/internal/model"
	"notification/pkg/httpx"
)

// registerRetryPolicyRoutes 注册重试策略相关路由。
func (s *Server) registerRetryPolicyRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/retry-policies", s.createRetryPolicy)
	mux.HandleFunc("GET /api/retry-policies", s.listRetryPolicies)
	mux.HandleFunc("GET /api/retry-policies/{id}", s.getRetryPolicy)
	mux.HandleFunc("PUT /api/retry-policies/{id}", s.updateRetryPolicy)
	mux.HandleFunc("DELETE /api/retry-policies/{id}", s.deleteRetryPolicy)
	mux.HandleFunc("POST /api/retry-policies/batch-delete", s.batchDeleteRetryPolicies)
}

type retryPolicyRequest struct {
	Name        string `json:"name"`
	ChannelType string `json:"channel_type"`
	MaxAttempts int    `json:"max_attempts"`
	BackoffMs   int64  `json:"backoff_ms"`
	Status      string `json:"status"`
}

func (s *Server) createRetryPolicy(w http.ResponseWriter, r *http.Request) {
	var req retryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.CreateRetryPolicy(model.RetryPolicy{
		Name: req.Name, ChannelType: req.ChannelType,
		MaxAttempts: req.MaxAttempts, BackoffMs: req.BackoffMs, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, p)
}

func (s *Server) listRetryPolicies(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RetryPolicyFilter{
		ChannelType: r.URL.Query().Get("channel_type"),
		Status:      r.URL.Query().Get("status"),
		Keyword:     r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRetryPolicies(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRetryPolicy(w http.ResponseWriter, r *http.Request) {
	p, err := s.svc.GetRetryPolicy(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) updateRetryPolicy(w http.ResponseWriter, r *http.Request) {
	var req retryPolicyRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	p, err := s.svc.UpdateRetryPolicy(r.PathValue("id"), model.RetryPolicy{
		Name: req.Name, ChannelType: req.ChannelType,
		MaxAttempts: req.MaxAttempts, BackoffMs: req.BackoffMs, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, p)
}

func (s *Server) deleteRetryPolicy(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteRetryPolicy(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) batchDeleteRetryPolicies(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 {
		httpx.BadRequest(w, "ids 不能为空")
		return
	}
	deleted, err := s.svc.BatchDeleteRetryPolicies(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"deleted": deleted})
}
