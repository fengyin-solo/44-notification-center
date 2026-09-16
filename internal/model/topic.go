package model

import (
	"strings"
	"time"
)

// 主题状态。
const (
	TopicActive   = "active"
	TopicInactive = "inactive"
)

// Topic 订阅主题。
type Topic struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	SubscriberCount int       `json:"subscriber_count"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// Validate 校验并规范化主题字段。
func (t *Topic) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.Description = strings.TrimSpace(t.Description)
	if t.Name == "" {
		return NewValidationError("name", "主题名称不能为空")
	}
	if t.SubscriberCount < 0 {
		return NewValidationError("subscriber_count", "订阅数不能为负")
	}
	if t.Status == "" {
		t.Status = TopicActive
	}
	if t.Status != TopicActive && t.Status != TopicInactive {
		return NewValidationError("status", "主题状态不合法")
	}
	return nil
}

// TopicFilter 主题筛选条件。
type TopicFilter struct {
	Status  string
	Keyword string
}

// Match 判断主题是否命中筛选条件。
func (f TopicFilter) Match(t *Topic) bool {
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) &&
			!strings.Contains(strings.ToLower(t.Description), k) {
			return false
		}
	}
	return true
}
