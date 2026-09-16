package model

import (
	"strings"
	"time"
)

// 消息优先级。
const (
	PriorityLow    = "low"
	PriorityNormal = "normal"
	PriorityHigh   = "high"
	PriorityUrgent = "urgent"
)

// 消息状态。
const (
	MessageDraft     = "draft"
	MessagePending   = "pending"
	MessageSending   = "sending"
	MessageSent      = "sent"
	MessageFailed    = "failed"
	MessageCancelled = "cancelled"
)

// messageTransitions 消息状态机。
var messageTransitions = map[string]map[string]bool{
	MessageDraft:     {MessagePending: true, MessageCancelled: true},
	MessagePending:   {MessageSending: true, MessageCancelled: true},
	MessageSending:   {MessageSent: true, MessageFailed: true},
	MessageSent:      {},
	MessageFailed:    {MessagePending: true},
	MessageCancelled: {},
}

// Message 待发送的通知消息。
type Message struct {
	ID          string     `json:"id"`
	TemplateID  string     `json:"template_id"`
	TopicID     string     `json:"topic_id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	ChannelType string     `json:"channel_type"`
	Priority    string     `json:"priority"`
	Status      string     `json:"status"`
	ScheduledAt *time.Time `json:"scheduled_at,omitempty"`
	SentAt      *time.Time `json:"sent_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Validate 校验并规范化消息字段。
func (m *Message) Validate() error {
	m.Title = strings.TrimSpace(m.Title)
	m.Content = strings.TrimSpace(m.Content)
	m.ChannelType = strings.TrimSpace(m.ChannelType)
	if m.Title == "" {
		return NewValidationError("title", "消息标题不能为空")
	}
	if m.Content == "" {
		return NewValidationError("content", "消息内容不能为空")
	}
	if !validChannelType(m.ChannelType) {
		return NewValidationError("channel_type", "渠道类型不合法")
	}
	if m.Priority == "" {
		m.Priority = PriorityNormal
	}
	if !validPriority(m.Priority) {
		return NewValidationError("priority", "消息优先级不合法")
	}
	if m.Status == "" {
		m.Status = MessageDraft
	}
	if !validMessageStatus(m.Status) {
		return NewValidationError("status", "消息状态不合法")
	}
	return nil
}

// validPriority 判断优先级是否合法。
func validPriority(p string) bool {
	switch p {
	case PriorityLow, PriorityNormal, PriorityHigh, PriorityUrgent:
		return true
	default:
		return false
	}
}

// validMessageStatus 判断消息状态是否合法。
func validMessageStatus(s string) bool {
	switch s {
	case MessageDraft, MessagePending, MessageSending, MessageSent, MessageFailed, MessageCancelled:
		return true
	default:
		return false
	}
}

// CanTransitionMessage 判断消息状态是否可流转。
func CanTransitionMessage(from, to string) bool {
	if m, ok := messageTransitions[from]; ok {
		return m[to]
	}
	return false
}

// MessageFilter 消息筛选条件。
type MessageFilter struct {
	TemplateID  string
	TopicID     string
	ChannelType string
	Priority    string
	Status      string
	Keyword     string
}

// Match 判断消息是否命中筛选条件。
func (f MessageFilter) Match(m *Message) bool {
	if f.TemplateID != "" && m.TemplateID != f.TemplateID {
		return false
	}
	if f.TopicID != "" && m.TopicID != f.TopicID {
		return false
	}
	if f.ChannelType != "" && m.ChannelType != f.ChannelType {
		return false
	}
	if f.Priority != "" && m.Priority != f.Priority {
		return false
	}
	if f.Status != "" && m.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(m.Title), k) &&
			!strings.Contains(strings.ToLower(m.Content), k) {
			return false
		}
	}
	return true
}
