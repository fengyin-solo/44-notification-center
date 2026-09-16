package handler

import (
	"net/http"

	"notification/internal/model"
	"notification/pkg/httpx"
)

// registerTemplateRoutes 注册模板相关路由。
func (s *Server) registerTemplateRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/templates", s.createTemplate)
	mux.HandleFunc("GET /api/templates", s.listTemplates)
	mux.HandleFunc("GET /api/templates/{id}", s.getTemplate)
	mux.HandleFunc("PUT /api/templates/{id}", s.updateTemplate)
	mux.HandleFunc("DELETE /api/templates/{id}", s.deleteTemplate)
	mux.HandleFunc("POST /api/templates/{id}/transition", s.transitionTemplate)
	mux.HandleFunc("POST /api/templates/{id}/render", s.renderTemplate)
}

type templateRequest struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Subject   string   `json:"subject"`
	Content   string   `json:"content"`
	Variables []string `json:"variables"`
	Status    string   `json:"status"`
}

func (s *Server) createTemplate(w http.ResponseWriter, r *http.Request) {
	var req templateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.CreateTemplate(model.Template{
		Name: req.Name, Type: req.Type, Subject: req.Subject,
		Content: req.Content, Variables: req.Variables, Status: req.Status,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, t)
}

func (s *Server) listTemplates(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.TemplateFilter{
		Type:    r.URL.Query().Get("type"),
		Status:  r.URL.Query().Get("status"),
		Keyword: r.URL.Query().Get("keyword"),
	}
	items, total, err := s.svc.ListTemplates(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getTemplate(w http.ResponseWriter, r *http.Request) {
	t, err := s.svc.GetTemplate(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) updateTemplate(w http.ResponseWriter, r *http.Request) {
	var req templateRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.UpdateTemplate(r.PathValue("id"), model.Template{
		Name: req.Name, Type: req.Type, Subject: req.Subject,
		Content: req.Content, Variables: req.Variables,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

func (s *Server) deleteTemplate(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteTemplate(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type transitionRequest struct {
	Status string `json:"status"`
}

func (s *Server) transitionTemplate(w http.ResponseWriter, r *http.Request) {
	var req transitionRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	t, err := s.svc.TransitionTemplate(r.PathValue("id"), req.Status)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, t)
}

type renderRequest struct {
	Variables map[string]string `json:"variables"`
}

func (s *Server) renderTemplate(w http.ResponseWriter, r *http.Request) {
	var req renderRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.RenderTemplate(r.PathValue("id"), req.Variables)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
