package handler

import (
	"net/http"

	"notification/internal/service"
	"notification/pkg/httpx"
)

// registerExportRoutes 注册数据导入导出相关路由。
func (s *Server) registerExportRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/export", s.exportSnapshot)
	mux.HandleFunc("POST /api/import", s.importSnapshot)
}

// exportSnapshot 导出全部数据的汇总快照 JSON。
func (s *Server) exportSnapshot(w http.ResponseWriter, r *http.Request) {
	snapshot, err := s.svc.ExportSnapshot()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, snapshot)
}

// importSnapshot 从快照 JSON 导入数据。
func (s *Server) importSnapshot(w http.ResponseWriter, r *http.Request) {
	var req service.ImportData
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	result, err := s.svc.ImportData(&req)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, result)
}
