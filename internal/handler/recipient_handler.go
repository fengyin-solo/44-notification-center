package handler

import (
	"net/http"

	"notification/internal/model"
	"notification/pkg/httpx"
)

// registerRecipientRoutes 注册接收人相关路由。
func (s *Server) registerRecipientRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/recipients", s.createRecipient)
	mux.HandleFunc("GET /api/recipients", s.listRecipients)
	mux.HandleFunc("GET /api/recipients/{id}", s.getRecipient)
	mux.HandleFunc("PUT /api/recipients/{id}", s.updateRecipient)
	mux.HandleFunc("DELETE /api/recipients/{id}", s.deleteRecipient)
	mux.HandleFunc("POST /api/recipients/{id}/status", s.setRecipientStatus)
	mux.HandleFunc("POST /api/recipients/batch-delete", s.batchDeleteRecipients)
}

type recipientRequest struct {
	Name        string `json:"name"`
	ChannelType string `json:"channel_type"`
	Address     string `json:"address"`
	Group       string `json:"group"`
	Status      string `json:"status"`
}

func (s *Server) createRecipient(w http.ResponseWriter, r *http.Request) {
	var req recipientRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.CreateRecipient(model.Recipient{
		Name: req.Name, ChannelType: req.ChannelType, Address: req.Address, Group: req.Group,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) listRecipients(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.RecipientFilter{
		ChannelType: r.URL.Query().Get("channel_type"),
		Group:       r.URL.Query().Get("group"),
		Status:      r.URL.Query().Get("status"),
		Keyword:     r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListRecipients(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getRecipient(w http.ResponseWriter, r *http.Request) {
	rec, err := s.svc.GetRecipient(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) updateRecipient(w http.ResponseWriter, r *http.Request) {
	var req recipientRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.UpdateRecipient(r.PathValue("id"), model.Recipient{
		Name: req.Name, ChannelType: req.ChannelType, Address: req.Address,
		Group: req.Group, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) deleteRecipient(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteRecipient(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

func (s *Server) setRecipientStatus(w http.ResponseWriter, r *http.Request) {
	var req transitionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.SetRecipientStatus(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) batchDeleteRecipients(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 {
		httpx.BadRequest(w, "ids 不能为空")
		return
	}
	deleted, err := s.svc.BatchDeleteRecipients(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"deleted": deleted})
}
