package handler

import (
	"net/http"
	"time"

	"notification/internal/model"
	"notification/pkg/httpx"
)

// registerMessageRoutes 注册消息相关路由。
func (s *Server) registerMessageRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/messages", s.createMessage)
	mux.HandleFunc("GET /api/messages", s.listMessages)
	mux.HandleFunc("GET /api/messages/{id}", s.getMessage)
	mux.HandleFunc("PUT /api/messages/{id}", s.updateMessage)
	mux.HandleFunc("DELETE /api/messages/{id}", s.deleteMessage)
	mux.HandleFunc("POST /api/messages/{id}/transition", s.transitionMessage)
	mux.HandleFunc("POST /api/messages/{id}/send", s.sendMessage)
	mux.HandleFunc("POST /api/messages/batch-send", s.batchSendMessages)
	mux.HandleFunc("POST /api/messages/batch-status", s.batchUpdateMessageStatus)
	mux.HandleFunc("POST /api/messages/batch-delete", s.batchDeleteMessages)
}

type messageRequest struct {
	TemplateID  string `json:"template_id"`
	TopicID     string `json:"topic_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	ChannelType string `json:"channel_type"`
	Priority    string `json:"priority"`
	Status      string `json:"status"`
	ScheduledAt string `json:"scheduled_at"`
}

func (s *Server) createMessage(w http.ResponseWriter, r *http.Request) {
	var req messageRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.CreateMessage(model.Message{
		TemplateID: req.TemplateID, TopicID: req.TopicID, Title: req.Title,
		Content: req.Content, ChannelType: req.ChannelType, Priority: req.Priority,
		Status: req.Status, ScheduledAt: parseOptionalTime(req.ScheduledAt),
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, m)
}

func (s *Server) listMessages(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.MessageFilter{
		TemplateID:  r.URL.Query().Get("template_id"),
		TopicID:     r.URL.Query().Get("topic_id"),
		ChannelType: r.URL.Query().Get("channel_type"),
		Priority:    r.URL.Query().Get("priority"),
		Status:      r.URL.Query().Get("status"),
		Keyword:     r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListMessages(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getMessage(w http.ResponseWriter, r *http.Request) {
	m, err := s.svc.GetMessage(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) updateMessage(w http.ResponseWriter, r *http.Request) {
	var req messageRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.UpdateMessage(r.PathValue("id"), model.Message{
		Title: req.Title, Content: req.Content, ChannelType: req.ChannelType, Priority: req.Priority,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) deleteMessage(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteMessage(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) transitionMessage(w http.ResponseWriter, r *http.Request) {
	var req transitionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	m, err := s.svc.TransitionMessage(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, m)
}

func (s *Server) sendMessage(w http.ResponseWriter, r *http.Request) {
	m, count, err := s.svc.SendMessage(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]interface{}{"message": m, "records": count})
}

func (s *Server) batchSendMessages(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 {
		httpx.BadRequest(w, "ids 不能为空")
		return
	}
	sent, err := s.svc.BatchSendMessages(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"sent": sent})
}

type batchStatusRequest struct {
	IDs    []string `json:"ids"`
	Status string   `json:"status"`
}

func (s *Server) batchUpdateMessageStatus(w http.ResponseWriter, r *http.Request) {
	var req batchStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 || req.Status == "" {
		httpx.BadRequest(w, "ids 与 status 不能为空")
		return
	}
	updated, err := s.svc.BatchUpdateMessageStatus(req.IDs, req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"updated": updated})
}

func (s *Server) batchDeleteMessages(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 {
		httpx.BadRequest(w, "ids 不能为空")
		return
	}
	deleted, err := s.svc.BatchDeleteMessages(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"deleted": deleted})
}

// parseOptionalTime 解析 RFC3339 时间，空串或非法返回 nil。
func parseOptionalTime(s string) *time.Time {
	if s == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil
	}
	return &t
}
