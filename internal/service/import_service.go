package service

import (
	"time"

	"notification/internal/model"
	"notification/pkg/idgen"
)

// ImportData 导入数据的载荷结构。
type ImportData struct {
	Channels      []*model.Channel      `json:"channels"`
	Templates     []*model.Template     `json:"templates"`
	Topics        []*model.Topic        `json:"topics"`
	Recipients    []*model.Recipient    `json:"recipients"`
	Subscriptions []*model.Subscription `json:"subscriptions"`
	Messages      []*model.Message      `json:"messages"`
	SendRecords   []*model.SendRecord   `json:"send_records"`
	RetryPolicies []*model.RetryPolicy  `json:"retry_policies"`
	Schedules     []*model.Schedule     `json:"schedules"`
}

// ImportResult 导入结果统计。
type ImportResult struct {
	Channels      int `json:"channels"`
	Templates     int `json:"templates"`
	Topics        int `json:"topics"`
	Recipients    int `json:"recipients"`
	Subscriptions int `json:"subscriptions"`
	Messages      int `json:"messages"`
	SendRecords   int `json:"send_records"`
	RetryPolicies int `json:"retry_policies"`
	Schedules     int `json:"schedules"`
	Total         int `json:"total"`
}

// ImportData 按实体类型导入数据，保留原 ID 并跳过冲突或非法记录。
func (s *Service) ImportData(data *ImportData) (*ImportResult, error) {
	if data == nil {
		return nil, model.NewValidationError("data", "导入数据不能为空")
	}
	result := &ImportResult{}
	now := time.Now()

	for _, c := range data.Channels {
		if c == nil {
			continue
		}
		if c.ID == "" {
			c.ID = idgen.Hex()
		}
		if c.CreatedAt.IsZero() {
			c.CreatedAt = now
		}
		if c.UpdatedAt.IsZero() {
			c.UpdatedAt = now
		}
		if err := c.Validate(); err != nil {
			continue
		}
		if s.store.CreateChannel(c) == nil {
			result.Channels++
		}
	}

	for _, t := range data.Templates {
		if t == nil {
			continue
		}
		if t.ID == "" {
			t.ID = idgen.Hex()
		}
		if t.CreatedAt.IsZero() {
			t.CreatedAt = now
		}
		if t.UpdatedAt.IsZero() {
			t.UpdatedAt = now
		}
		if err := t.Validate(); err != nil {
			continue
		}
		if s.store.CreateTemplate(t) == nil {
			result.Templates++
		}
	}

	for _, t := range data.Topics {
		if t == nil {
			continue
		}
		if t.ID == "" {
			t.ID = idgen.Hex()
		}
		if t.CreatedAt.IsZero() {
			t.CreatedAt = now
		}
		if t.UpdatedAt.IsZero() {
			t.UpdatedAt = now
		}
		if err := t.Validate(); err != nil {
			continue
		}
		if s.store.CreateTopic(t) == nil {
			result.Topics++
		}
	}

	for _, r := range data.Recipients {
		if r == nil {
			continue
		}
		if r.ID == "" {
			r.ID = idgen.Hex()
		}
		if r.CreatedAt.IsZero() {
			r.CreatedAt = now
		}
		if r.UpdatedAt.IsZero() {
			r.UpdatedAt = now
		}
		if err := r.Validate(); err != nil {
			continue
		}
		if s.store.CreateRecipient(r) == nil {
			result.Recipients++
		}
	}

	for _, sub := range data.Subscriptions {
		if sub == nil {
			continue
		}
		if sub.ID == "" {
			sub.ID = idgen.Hex()
		}
		if sub.CreatedAt.IsZero() {
			sub.CreatedAt = now
		}
		if sub.UpdatedAt.IsZero() {
			sub.UpdatedAt = now
		}
		if err := sub.Validate(); err != nil {
			continue
		}
		if _, err := s.store.GetTopic(sub.TopicID); err != nil {
			continue
		}
		if _, err := s.store.GetRecipient(sub.RecipientID); err != nil {
			continue
		}
		if s.store.CreateSubscription(sub) == nil {
			result.Subscriptions++
		}
	}

	for _, m := range data.Messages {
		if m == nil {
			continue
		}
		if m.ID == "" {
			m.ID = idgen.Hex()
		}
		if m.CreatedAt.IsZero() {
			m.CreatedAt = now
		}
		if m.UpdatedAt.IsZero() {
			m.UpdatedAt = now
		}
		if err := m.Validate(); err != nil {
			continue
		}
		if s.store.CreateMessage(m) == nil {
			result.Messages++
		}
	}

	for _, r := range data.SendRecords {
		if r == nil {
			continue
		}
		if r.ID == "" {
			r.ID = idgen.Hex()
		}
		if r.CreatedAt.IsZero() {
			r.CreatedAt = now
		}
		if err := r.Validate(); err != nil {
			continue
		}
		if s.store.CreateSendRecord(r) == nil {
			result.SendRecords++
		}
	}

	for _, p := range data.RetryPolicies {
		if p == nil {
			continue
		}
		if p.ID == "" {
			p.ID = idgen.Hex()
		}
		if p.CreatedAt.IsZero() {
			p.CreatedAt = now
		}
		if p.UpdatedAt.IsZero() {
			p.UpdatedAt = now
		}
		if err := p.Validate(); err != nil {
			continue
		}
		if s.store.CreateRetryPolicy(p) == nil {
			result.RetryPolicies++
		}
	}

	for _, sc := range data.Schedules {
		if sc == nil {
			continue
		}
		if sc.ID == "" {
			sc.ID = idgen.Hex()
		}
		if sc.CreatedAt.IsZero() {
			sc.CreatedAt = now
		}
		if sc.UpdatedAt.IsZero() {
			sc.UpdatedAt = now
		}
		if err := sc.Validate(); err != nil {
			continue
		}
		if s.store.CreateSchedule(sc) == nil {
			result.Schedules++
		}
	}

	result.Total = result.Channels + result.Templates + result.Topics + result.Recipients +
		result.Subscriptions + result.Messages + result.SendRecords + result.RetryPolicies + result.Schedules
	s.log.Infof("数据导入完成: total=%d", result.Total)
	return result, nil
}
