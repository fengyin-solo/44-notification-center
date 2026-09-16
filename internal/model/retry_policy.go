package model

import (
	"strings"
	"time"
)

// 重试策略状态。
const (
	RetryPolicyEnabled  = "enabled"
	RetryPolicyDisabled = "disabled"
)

// RetryPolicy 失败重试策略。
type RetryPolicy struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	ChannelType string    `json:"channel_type"`
	MaxAttempts int       `json:"max_attempts"`
	BackoffMs   int64     `json:"backoff_ms"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验并规范化重试策略字段。
func (p *RetryPolicy) Validate() error {
	p.Name = strings.TrimSpace(p.Name)
	p.ChannelType = strings.TrimSpace(p.ChannelType)
	if p.Name == "" {
		return NewValidationError("name", "策略名称不能为空")
	}
	if !validChannelType(p.ChannelType) {
		return NewValidationError("channel_type", "渠道类型不合法")
	}
	if p.MaxAttempts <= 0 {
		return NewValidationError("max_attempts", "最大尝试次数必须大于 0")
	}
	if p.BackoffMs < 0 {
		return NewValidationError("backoff_ms", "退避间隔不能为负")
	}
	if p.Status == "" {
		p.Status = RetryPolicyEnabled
	}
	if p.Status != RetryPolicyEnabled && p.Status != RetryPolicyDisabled {
		return NewValidationError("status", "策略状态不合法")
	}
	return nil
}

// RetryPolicyFilter 重试策略筛选条件。
type RetryPolicyFilter struct {
	ChannelType string
	Status      string
	Keyword     string
}

// Match 判断重试策略是否命中筛选条件。
func (f RetryPolicyFilter) Match(p *RetryPolicy) bool {
	if f.ChannelType != "" && p.ChannelType != f.ChannelType {
		return false
	}
	if f.Status != "" && p.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(p.Name), k) {
			return false
		}
	}
	return true
}
