// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"notification/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

// Store 聚合全部实体的数据访问方法，便于测试时替换实现。
type Store interface {
	// Channel 渠道。
	CreateChannel(c *model.Channel) error
	GetChannel(id string) (*model.Channel, error)
	GetChannelByName(name string) (*model.Channel, error)
	ListChannels() []*model.Channel
	UpdateChannel(c *model.Channel) error
	DeleteChannel(id string) error

	// Template 模板。
	CreateTemplate(t *model.Template) error
	GetTemplate(id string) (*model.Template, error)
	GetTemplateByName(name string) (*model.Template, error)
	ListTemplates() []*model.Template
	UpdateTemplate(t *model.Template) error
	DeleteTemplate(id string) error

	// Topic 主题。
	CreateTopic(t *model.Topic) error
	GetTopic(id string) (*model.Topic, error)
	ListTopics() []*model.Topic
	UpdateTopic(t *model.Topic) error
	DeleteTopic(id string) error

	// Recipient 接收人。
	CreateRecipient(r *model.Recipient) error
	GetRecipient(id string) (*model.Recipient, error)
	GetRecipientByAddress(channelType, address string) (*model.Recipient, error)
	ListRecipients() []*model.Recipient
	UpdateRecipient(r *model.Recipient) error
	DeleteRecipient(id string) error

	// Subscription 订阅。
	CreateSubscription(s *model.Subscription) error
	GetSubscription(id string) (*model.Subscription, error)
	GetSubscriptionByTopicRecipient(topicID, recipientID string) (*model.Subscription, error)
	ListSubscriptions() []*model.Subscription
	UpdateSubscription(s *model.Subscription) error
	DeleteSubscription(id string) error

	// Message 消息。
	CreateMessage(m *model.Message) error
	GetMessage(id string) (*model.Message, error)
	ListMessages() []*model.Message
	UpdateMessage(m *model.Message) error
	DeleteMessage(id string) error

	// SendRecord 发送记录。
	CreateSendRecord(r *model.SendRecord) error
	GetSendRecord(id string) (*model.SendRecord, error)
	ListSendRecords() []*model.SendRecord
	UpdateSendRecord(r *model.SendRecord) error
	DeleteSendRecord(id string) error

	// RetryPolicy 重试策略。
	CreateRetryPolicy(p *model.RetryPolicy) error
	GetRetryPolicy(id string) (*model.RetryPolicy, error)
	GetRetryPolicyByName(name string) (*model.RetryPolicy, error)
	ListRetryPolicies() []*model.RetryPolicy
	UpdateRetryPolicy(p *model.RetryPolicy) error
	DeleteRetryPolicy(id string) error

	// Schedule 定时任务。
	CreateSchedule(s *model.Schedule) error
	GetSchedule(id string) (*model.Schedule, error)
	ListSchedules() []*model.Schedule
	UpdateSchedule(s *model.Schedule) error
	DeleteSchedule(id string) error
}
