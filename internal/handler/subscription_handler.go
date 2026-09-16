package handler

import (
	"net/http"

	"notification/internal/model"
	"notification/pkg/httpx"
)

// registerSubscriptionRoutes 注册订阅相关路由。
func (s *Server) registerSubscriptionRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/subscriptions", s.createSubscription)
	mux.HandleFunc("GET /api/subscriptions", s.listSubscriptions)
	mux.HandleFunc("GET /api/subscriptions/{id}", s.getSubscription)
	mux.HandleFunc("PUT /api/subscriptions/{id}", s.updateSubscription)
	mux.HandleFunc("DELETE /api/subscriptions/{id}", s.deleteSubscription)
	mux.HandleFunc("POST /api/subscriptions/{id}/unsubscribe", s.unsubscribe)
	mux.HandleFunc("POST /api/subscriptions/{id}/subscribe", s.subscribe)
	mux.HandleFunc("POST /api/subscriptions/batch-delete", s.batchDeleteSubscriptions)
}

type subscriptionRequest struct {
	TopicID     string `json:"topic_id"`
	RecipientID string `json:"recipient_id"`
	ChannelType string `json:"channel_type"`
	Status      string `json:"status"`
}

func (s *Server) createSubscription(w http.ResponseWriter, r *http.Request) {
	var req subscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.CreateSubscription(model.Subscription{
		TopicID: req.TopicID, RecipientID: req.RecipientID, ChannelType: req.ChannelType,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, sub)
}

func (s *Server) listSubscriptions(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SubscriptionFilter{
		TopicID:     r.URL.Query().Get("topic_id"),
		RecipientID: r.URL.Query().Get("recipient_id"),
		ChannelType: r.URL.Query().Get("channel_type"),
		Status:      r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSubscriptions(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSubscription(w http.ResponseWriter, r *http.Request) {
	sub, err := s.svc.GetSubscription(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) updateSubscription(w http.ResponseWriter, r *http.Request) {
	var req subscriptionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	sub, err := s.svc.UpdateSubscription(r.PathValue("id"), model.Subscription{
		ChannelType: req.ChannelType, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) deleteSubscription(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSubscription(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) unsubscribe(w http.ResponseWriter, r *http.Request) {
	sub, err := s.svc.Unsubscribe(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) subscribe(w http.ResponseWriter, r *http.Request) {
	sub, err := s.svc.Subscribe(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, sub)
}

func (s *Server) batchDeleteSubscriptions(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 {
		httpx.BadRequest(w, "ids 不能为空")
		return
	}
	deleted, err := s.svc.BatchDeleteSubscriptions(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"deleted": deleted})
}
