package handler

import (
	"net/http"

	"notification/internal/model"
	"notification/pkg/httpx"
)

// registerChannelRoutes 注册渠道相关路由。
func (s *Server) registerChannelRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/channels", s.createChannel)
	mux.HandleFunc("GET /api/channels", s.listChannels)
	mux.HandleFunc("GET /api/channels/{id}", s.getChannel)
	mux.HandleFunc("PUT /api/channels/{id}", s.updateChannel)
	mux.HandleFunc("DELETE /api/channels/{id}", s.deleteChannel)
	mux.HandleFunc("POST /api/channels/batch-delete", s.batchDeleteChannels)
}

type channelRequest struct {
	Name     string `json:"name"`
	Type     string `json:"type"`
	Status   string `json:"status"`
	Config   string `json:"config"`
	Priority int    `json:"priority"`
}

func (s *Server) createChannel(w http.ResponseWriter, r *http.Request) {
	var req channelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.CreateChannel(model.Channel{
		Name: req.Name, Type: req.Type, Status: req.Status, Config: req.Config, Priority: req.Priority,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, c)
}

func (s *Server) listChannels(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.ChannelFilter{
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListChannels(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getChannel(w http.ResponseWriter, r *http.Request) {
	c, err := s.svc.GetChannel(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) updateChannel(w http.ResponseWriter, r *http.Request) {
	var req channelRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	c, err := s.svc.UpdateChannel(r.PathValue("id"), model.Channel{
		Name: req.Name, Type: req.Type, Status: req.Status, Config: req.Config, Priority: req.Priority,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, c)
}

func (s *Server) deleteChannel(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteChannel(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type batchDeleteRequest struct {
	IDs []string `json:"ids"`
}

func (s *Server) batchDeleteChannels(w http.ResponseWriter, r *http.Request) {
	var req batchDeleteRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 {
		httpx.BadRequest(w, "ids 不能为空")
		return
	}
	deleted, err := s.svc.BatchDeleteChannels(req.IDs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"deleted": deleted})
}
