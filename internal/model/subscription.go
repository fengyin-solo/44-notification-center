package model

import (
	"time"
)

// 订阅状态。
const (
	SubscriptionSubscribed   = "subscribed"
	SubscriptionUnsubscribed = "unsubscribed"
)

// subscriptionTransitions 订阅状态机：subscribed→unsubscribed。
var subscriptionTransitions = map[string]map[string]bool{
	SubscriptionSubscribed:   {SubscriptionUnsubscribed: true},
	SubscriptionUnsubscribed: {},
}

// Subscription 主题订阅关系。
type Subscription struct {
	ID          string    `json:"id"`
	TopicID     string    `json:"topic_id"`
	RecipientID string    `json:"recipient_id"`
	ChannelType string    `json:"channel_type"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Validate 校验订阅字段（跨实体存在性在 service 层校验）。
func (s *Subscription) Validate() error {
	if s.TopicID == "" {
		return NewValidationError("topic_id", "主题 ID 不能为空")
	}
	if s.RecipientID == "" {
		return NewValidationError("recipient_id", "接收人 ID 不能为空")
	}
	if !validChannelType(s.ChannelType) {
		return NewValidationError("channel_type", "渠道类型不合法")
	}
	if s.Status == "" {
		s.Status = SubscriptionSubscribed
	}
	if s.Status != SubscriptionSubscribed && s.Status != SubscriptionUnsubscribed {
		return NewValidationError("status", "订阅状态不合法")
	}
	return nil
}

// CanTransitionSubscription 判断订阅状态是否可流转。
func CanTransitionSubscription(from, to string) bool {
	if m, ok := subscriptionTransitions[from]; ok {
		return m[to]
	}
	return false
}

// SubscriptionFilter 订阅筛选条件。
type SubscriptionFilter struct {
	TopicID     string
	RecipientID string
	ChannelType string
	Status      string
}

// Match 判断订阅是否命中筛选条件。
func (f SubscriptionFilter) Match(s *Subscription) bool {
	if f.TopicID != "" && s.TopicID != f.TopicID {
		return false
	}
	if f.RecipientID != "" && s.RecipientID != f.RecipientID {
		return false
	}
	if f.ChannelType != "" && s.ChannelType != f.ChannelType {
		return false
	}
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	return true
}
