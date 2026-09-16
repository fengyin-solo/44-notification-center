package handler

import (
	"net/http"
	"strconv"

	"notification/pkg/httpx"
)

// registerStatsRoutes 注册统计报表相关路由。
func (s *Server) registerStatsRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/stats/overview", s.statsOverview)
	mux.HandleFunc("GET /api/stats/channels", s.statsChannels)
	mux.HandleFunc("GET /api/stats/message-status", s.statsMessageStatus)
	mux.HandleFunc("GET /api/stats/message-priority", s.statsMessagePriority)
	mux.HandleFunc("GET /api/stats/top-templates", s.statsTopTemplates)
	mux.HandleFunc("GET /api/stats/topic-subscriptions", s.statsTopicSubscriptions)
}

func (s *Server) statsOverview(w http.ResponseWriter, r *http.Request) {
	o, err := s.svc.Overview()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, o)
}

func (s *Server) statsChannels(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.ChannelBreakdown()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) statsMessageStatus(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.MessageStatusBreakdown()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) statsMessagePriority(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.MessagePriorityBreakdown()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) statsTopTemplates(w http.ResponseWriter, r *http.Request) {
	n := 5
	if v, err := strconv.Atoi(r.URL.Query().Get("n")); err == nil && v > 0 {
		n = v
	}
	list, err := s.svc.TopTemplates(n)
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}

func (s *Server) statsTopicSubscriptions(w http.ResponseWriter, r *http.Request) {
	list, err := s.svc.TopicSubscriptionBreakdown()
	if err != nil {
		writeServiceError(w, err)
		return
	}
	httpx.OK(w, list)
}
