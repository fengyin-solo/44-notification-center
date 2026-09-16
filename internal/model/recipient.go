package model

import (
	"strings"
	"time"
)

// 接收人状态。
const (
	RecipientActive       = "active"
	RecipientUnsubscribed = "unsubscribed"
)

// Recipient 通知接收人。
type Recipient struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ChannelType string    `json:"channel_type"`
	Address     string    `json:"address"`
	Group       string    `json:"group"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验并规范化接收人字段。
func (r *Recipient) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.ChannelType = strings.TrimSpace(r.ChannelType)
	r.Address = strings.TrimSpace(r.Address)
	r.Group = strings.TrimSpace(r.Group)
	if r.Name == "" {
		return NewValidationError("name", "接收人姓名不能为空")
	}
	if !validChannelType(r.ChannelType) {
		return NewValidationError("channel_type", "接收渠道类型不合法")
	}
	if r.Address == "" {
		return NewValidationError("address", "接收地址不能为空")
	}
	if r.Status == "" {
		r.Status = RecipientActive
	}
	if r.Status != RecipientActive && r.Status != RecipientUnsubscribed {
		return NewValidationError("status", "接收人状态不合法")
	}
	return nil
}

// RecipientFilter 接收人筛选条件。
type RecipientFilter struct {
	ChannelType string
	Group       string
	Status      string
	Keyword     string
}

// Match 判断接收人是否命中筛选条件。
func (f RecipientFilter) Match(r *Recipient) bool {
	if f.ChannelType != "" && r.ChannelType != f.ChannelType {
		return false
	}
	if f.Group != "" && r.Group != f.Group {
		return false
	}
	if f.Status != "" && r.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(r.Name), k) &&
			!strings.Contains(strings.ToLower(r.Address), k) {
			return false
		}
	}
	return true
}
