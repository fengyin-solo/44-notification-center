package model

import (
	"time"
)

// 发送记录状态。
const (
	SendRecordPending  = "pending"
	SendRecordSuccess  = "success"
	SendRecordFailed   = "failed"
	SendRecordRetrying = "retrying"
)

// sendRecordTransitions 发送记录状态机。
var sendRecordTransitions = map[string]map[string]bool{
	SendRecordPending:  {SendRecordSuccess: true, SendRecordFailed: true},
	SendRecordRetrying: {SendRecordSuccess: true, SendRecordFailed: true},
	SendRecordFailed:   {SendRecordRetrying: true},
	SendRecordSuccess:  {},
}

// SendRecord 单条消息的发送记录。
type SendRecord struct {
	ID          string     `json:"id"`
	MessageID   string     `json:"message_id"`
	RecipientID string     `json:"recipient_id"`
	ChannelType string     `json:"channel_type"`
	Status      string     `json:"status"`
	Attempts    int        `json:"attempts"`
	DurationMs  int64      `json:"duration_ms"`
	Error       string     `json:"error,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

// Validate 校验发送记录字段。
func (r *SendRecord) Validate() error {
	if r.MessageID == "" {
		return NewValidationError("message_id", "消息 ID 不能为空")
	}
	if r.RecipientID == "" {
		return NewValidationError("recipient_id", "接收人 ID 不能为空")
	}
	if !validChannelType(r.ChannelType) {
		return NewValidationError("channel_type", "渠道类型不合法")
	}
	if r.Status == "" {
		r.Status = SendRecordPending
	}
	if !validSendRecordStatus(r.Status) {
		return NewValidationError("status", "发送记录状态不合法")
	}
	if r.Attempts < 0 {
		return NewValidationError("attempts", "尝试次数不能为负")
	}
	if r.DurationMs < 0 {
		return NewValidationError("duration_ms", "耗时不能为负")
	}
	return nil
}

// validSendRecordStatus 判断发送记录状态是否合法。
func validSendRecordStatus(s string) bool {
	switch s {
	case SendRecordPending, SendRecordSuccess, SendRecordFailed, SendRecordRetrying:
		return true
	default:
		return false
	}
}

// CanTransitionSendRecord 判断发送记录状态是否可流转。
func CanTransitionSendRecord(from, to string) bool {
	if m, ok := sendRecordTransitions[from]; ok {
		return m[to]
	}
	return false
}

// SendRecordFilter 发送记录筛选条件。
type SendRecordFilter struct {
	MessageID   string
	RecipientID string
	ChannelType string
	Status      string
}

// Match 判断发送记录是否命中筛选条件。
func (f SendRecordFilter) Match(r *SendRecord) bool {
	if f.MessageID != "" && r.MessageID != f.MessageID {
		return false
	}
	if f.RecipientID != "" && r.RecipientID != f.RecipientID {
		return false
	}
	if f.ChannelType != "" && r.ChannelType != f.ChannelType {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	return true
}
