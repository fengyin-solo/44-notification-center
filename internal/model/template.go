package model

import (
	"strings"
	"time"
)

// 模板状态。
const (
	TemplateDraft    = "draft"
	TemplateActive   = "active"
	TemplateArchived = "archived"
)

// templateTransitions 模板状态机：draft→active→archived。
var templateTransitions = map[string]map[string]bool{
	TemplateDraft:    {TemplateActive: true, TemplateArchived: true},
	TemplateActive:   {TemplateArchived: true, TemplateDraft: true},
	TemplateArchived: {},
}

// Template 消息模板。
type Template struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Subject   string    `json:"subject"`
	Content   string    `json:"content"`
	Variables []string  `json:"variables"`
	Version   int       `json:"version"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Validate 校验并规范化模板字段。
func (t *Template) Validate() error {
	t.Name = strings.TrimSpace(t.Name)
	t.Type = strings.TrimSpace(t.Type)
	t.Subject = strings.TrimSpace(t.Subject)
	t.Content = strings.TrimSpace(t.Content)
	if t.Name == "" {
		return NewValidationError("name", "模板名称不能为空")
	}
	if !validChannelType(t.Type) {
		return NewValidationError("type", "模板类型不合法")
	}
	if t.Content == "" {
		return NewValidationError("content", "模板内容不能为空")
	}
	if t.Version <= 0 {
		t.Version = 1
	}
	if t.Status == "" {
		t.Status = TemplateDraft
	}
	if !validTemplateStatus(t.Status) {
		return NewValidationError("status", "模板状态不合法")
	}
	return nil
}

// validTemplateStatus 判断模板状态是否合法。
func validTemplateStatus(s string) bool {
	switch s {
	case TemplateDraft, TemplateActive, TemplateArchived:
		return true
	default:
		return false
	}
}

// CanTransition 判断模板状态是否可流转。
func CanTransitionTemplate(from, to string) bool {
	if m, ok := templateTransitions[from]; ok {
		return m[to]
	}
	return false
}

// TemplateFilter 模板筛选条件。
type TemplateFilter struct {
	Type    string
	Status  string
	Keyword string
}

// Match 判断模板是否命中筛选条件。
func (f TemplateFilter) Match(t *Template) bool {
	if f.Type != "" && t.Type != f.Type {
		return false
	}
	if f.Status != "" && t.Status != f.Status {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(t.Name), k) &&
			!strings.Contains(strings.ToLower(t.Subject), k) {
			return false
		}
	}
	return true
}
