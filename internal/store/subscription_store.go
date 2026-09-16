package store

import "notification/internal/model"

// CreateSubscription 创建订阅，同一接收人对同一主题只能有一条订阅。
func (s *MemoryStore) CreateSubscription(sub *model.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.subscriptions {
		if exist.TopicID == sub.TopicID && exist.RecipientID == sub.RecipientID {
			return ErrConflict
		}
	}
	s.subscriptions[sub.ID] = sub
	return nil
}

// GetSubscription 按 ID 获取订阅。
func (s *MemoryStore) GetSubscription(id string) (*model.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sub, ok := s.subscriptions[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sub, nil
}

// GetSubscriptionByTopicRecipient 按主题+接收人获取订阅。
func (s *MemoryStore) GetSubscriptionByTopicRecipient(topicID, recipientID string) (*model.Subscription, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, sub := range s.subscriptions {
		if sub.TopicID == topicID && sub.RecipientID == recipientID {
			return sub, nil
		}
	}
	return nil, ErrNotFound
}

// ListSubscriptions 返回全部订阅。
func (s *MemoryStore) ListSubscriptions() []*model.Subscription {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Subscription, 0, len(s.subscriptions))
	for _, sub := range s.subscriptions {
		list = append(list, sub)
	}
	return list
}

// UpdateSubscription 更新订阅。
func (s *MemoryStore) UpdateSubscription(sub *model.Subscription) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscriptions[sub.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.subscriptions {
		if exist.ID != sub.ID && exist.TopicID == sub.TopicID && exist.RecipientID == sub.RecipientID {
			return ErrConflict
		}
	}
	s.subscriptions[sub.ID] = sub
	return nil
}

// DeleteSubscription 删除订阅。
func (s *MemoryStore) DeleteSubscription(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.subscriptions[id]; !ok {
		return ErrNotFound
	}
	delete(s.subscriptions, id)
	return nil
}
