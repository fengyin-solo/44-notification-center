package model

import (
	"strings"
	"time"
)

// 定时任务状态。
const (
	SchedulePending   = "pending"
	ScheduleExecuted  = "executed"
	ScheduleCancelled = "cancelled"
)

// scheduleTransitions 定时任务状态机：pending→executed/cancelled。
var scheduleTransitions = map[string]map[string]bool{
	SchedulePending:   {ScheduleExecuted: true, ScheduleCancelled: true},
	ScheduleExecuted:  {},
	ScheduleCancelled: {},
}

// Schedule 消息定时发送任务。
type Schedule struct {
	ID        string     `json:"id"`
	MessageID string     `json:"message_id"`
	CronExpr  string     `json:"cron_expr"`
	Status    string     `json:"status"`
	NextRunAt *time.Time `json:"next_run_at,omitempty"`
	LastRunAt *time.Time `json:"last_run_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// Validate 校验并规范化定时任务字段。
func (s *Schedule) Validate() error {
	s.MessageID = strings.TrimSpace(s.MessageID)
	s.CronExpr = strings.TrimSpace(s.CronExpr)
	if s.MessageID == "" {
		return NewValidationError("message_id", "消息 ID 不能为空")
	}
	if s.CronExpr == "" {
		return NewValidationError("cron_expr", "Cron 表达式不能为空")
	}
	if err := ParseCron(s.CronExpr); err != nil {
		return err
	}
	if s.Status == "" {
		s.Status = SchedulePending
	}
	if !validScheduleStatus(s.Status) {
		return NewValidationError("status", "定时任务状态不合法")
	}
	return nil
}

// validScheduleStatus 判断定时任务状态是否合法。
func validScheduleStatus(s string) bool {
	switch s {
	case SchedulePending, ScheduleExecuted, ScheduleCancelled:
		return true
	default:
		return false
	}
}

// CanTransitionSchedule 判断定时任务状态是否可流转。
func CanTransitionSchedule(from, to string) bool {
	if m, ok := scheduleTransitions[from]; ok {
		return m[to]
	}
	return false
}

// ScheduleFilter 定时任务筛选条件。
type ScheduleFilter struct {
	MessageID string
	Status    string
}

// Match 判断定时任务是否命中筛选条件。
func (f ScheduleFilter) Match(s *Schedule) bool {
	if f.MessageID != "" && s.MessageID != f.MessageID {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}
