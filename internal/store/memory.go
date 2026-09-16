package store

import (
	"sync"

	"notification/internal/model"
)

// MemoryStore 基于内存的 Store 实现，线程安全。
type MemoryStore struct {
	mu            sync.RWMutex
	channels      map[string]*model.Channel
	templates     map[string]*model.Template
	topics        map[string]*model.Topic
	recipients    map[string]*model.Recipient
	subscriptions map[string]*model.Subscription
	messages      map[string]*model.Message
	sendRecords   map[string]*model.SendRecord
	retryPolicies map[string]*model.RetryPolicy
	schedules     map[string]*model.Schedule
}

// NewMemoryStore 创建空的内存存储。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		channels:      make(map[string]*model.Channel),
		templates:     make(map[string]*model.Template),
		topics:        make(map[string]*model.Topic),
		recipients:    make(map[string]*model.Recipient),
		subscriptions: make(map[string]*model.Subscription),
		messages:      make(map[string]*model.Message),
		sendRecords:   make(map[string]*model.SendRecord),
		retryPolicies: make(map[string]*model.RetryPolicy),
		schedules:     make(map[string]*model.Schedule),
	}
}

var _ Store = (*MemoryStore)(nil)
