package model

import (
	"strings"
	"time"
)

// 渠道类型。
const (
	ChannelTypeEmail   = "email"
	ChannelTypeSMS     = "sms"
	ChannelTypePush    = "push"
	ChannelTypeWebhook = "webhook"
)

// 渠道状态。
const (
	ChannelEnabled  = "enabled"
	ChannelDisabled = "disabled"
)

// Channel 通知渠道。
type Channel struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Config    string    `json:"config"`
	Priority  int       `json:"priority"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化渠道字段。
func (c *Channel) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Type = strings.TrimSpace(c.Type)
	c.Config = strings.TrimSpace(c.Config)
	if c.Name == "" {
		return NewValidationError("name", "渠道名称不能为空")
	}
	if !validChannelType(c.Type) {
		return NewValidationError("type", "渠道类型不合法")
	}
	if c.Status == "" {
		c.Status = ChannelEnabled
	}
	if c.Status != ChannelEnabled && c.Status != ChannelDisabled {
		return NewValidationError("status", "渠道状态不合法")
	}
	if c.Priority < 0 {
		return NewValidationError("priority", "渠道优先级不能为负")
	}
	return nil
}

// validChannelType 判断渠道类型是否合法。
func validChannelType(t string) bool {
	switch t {
	case ChannelTypeEmail, ChannelTypeSMS, ChannelTypePush, ChannelTypeWebhook:
		return true
	default:
		return false
	}
}

// ChannelFilter 渠道筛选条件。
type ChannelFilter struct {
	Type    string
	Status  string
	Keyword string
}

// Match 判断渠道是否命中筛选条件。
func (f ChannelFilter) Match(c *Channel) bool {
	if f.Type != "" && c.Type != f.Type {
		return false
	}
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) {
			return false
		}
	}
	return true
}
