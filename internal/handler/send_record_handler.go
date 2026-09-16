package handler

import (
	"net/http"

	"notification/internal/model"
	"notification/pkg/httpx"
)

// registerSendRecordRoutes 注册发送记录相关路由。
func (s *Server) registerSendRecordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/send-records", s.createSendRecord)
	mux.HandleFunc("GET /api/send-records", s.listSendRecords)
	mux.HandleFunc("GET /api/send-records/{id}", s.getSendRecord)
	mux.HandleFunc("PUT /api/send-records/{id}", s.updateSendRecord)
	mux.HandleFunc("DELETE /api/send-records/{id}", s.deleteSendRecord)
	mux.HandleFunc("POST /api/send-records/{id}/success", s.markRecordSuccess)
	mux.HandleFunc("POST /api/send-records/{id}/failed", s.markRecordFailed)
	mux.HandleFunc("POST /api/send-records/{id}/retry", s.retryRecord)
	mux.HandleFunc("POST /api/send-records/batch-status", s.batchUpdateRecordStatus)
}

type sendRecordRequest struct {
	MessageID   string `json:"message_id"`
	RecipientID string `json:"recipient_id"`
	ChannelType string `json:"channel_type"`
	Status      string `json:"status"`
	Attempts    int    `json:"attempts"`
	DurationMs  int64  `json:"duration_ms"`
	Error       string `json:"error"`
}

func (s *Server) createSendRecord(w http.ResponseWriter, r *http.Request) {
	var req sendRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.CreateSendRecord(model.SendRecord{
		MessageID: req.MessageID, RecipientID: req.RecipientID, ChannelType: req.ChannelType,
		Status: req.Status, Attempts: req.Attempts, DurationMs: req.DurationMs, Error: req.Error,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.Created(w, rec)
}

func (s *Server) listSendRecords(w http.ResponseWriter, r *http.Request) {
	pp := httpx.ParsePagination(r, 20, s.maxPageSize())
	filter := model.SendRecordFilter{
		MessageID:   r.URL.Query().Get("message_id"),
		RecipientID: r.URL.Query().Get("recipient_id"),
		ChannelType: r.URL.Query().Get("channel_type"),
		Status:      r.URL.Query().Get("status"),
	}
	items, total, err := s.svc.ListSendRecords(filter, pp.Page, pp.Size)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, httpx.PageResult{
		Items:      items,
		Pagination: httpx.Pagination{Page: pp.Page, Size: pp.Size, Total: total},
	})
}

func (s *Server) getSendRecord(w http.ResponseWriter, r *http.Request) {
	rec, err := s.svc.GetSendRecord(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) updateSendRecord(w http.ResponseWriter, r *http.Request) {
	var req sendRecordRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.UpdateSendRecord(r.PathValue("id"), model.SendRecord{
		Status: req.Status, Attempts: req.Attempts, DurationMs: req.DurationMs, Error: req.Error,
	})
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) deleteSendRecord(w http.ResponseWriter, r *http.Request) {
	if err := s.svc.DeleteSendRecord(r.PathValue("id")); err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.NoContent(w)
}

type recordSuccessRequest struct {
	DurationMs int64 `json:"duration_ms"`
}

func (s *Server) markRecordSuccess(w http.ResponseWriter, r *http.Request) {
	var req recordSuccessRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.MarkRecordSuccess(r.PathValue("id"), req.DurationMs)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

type recordFailedRequest struct {
	Error string `json:"error"`
}

func (s *Server) markRecordFailed(w http.ResponseWriter, r *http.Request) {
	var req recordFailedRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	rec, err := s.svc.MarkRecordFailed(r.PathValue("id"), req.Error)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) retryRecord(w http.ResponseWriter, r *http.Request) {
	rec, err := s.svc.RetryRecord(r.PathValue("id"))
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, rec)
}

func (s *Server) batchUpdateRecordStatus(w http.ResponseWriter, r *http.Request) {
	var req batchStatusRequest
	if err := httpx.Decode(r, &req); err != nil {
		httpx.BadRequest(w, "请求体解析失败: "+err.Error())
		return
	}
	if len(req.IDs) == 0 || req.Status == "" {
		httpx.BadRequest(w, "ids 与 status 不能为空")
		return
	}
	updated, err := s.svc.BatchUpdateRecordStatus(req.IDs, req.Status, "")
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, map[string]int{"updated": updated})
}
