package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/idgen"
)

// CreateSubscription 创建订阅，校验主题与接收人均存在。
func (s *Service) CreateSubscription(input model.Subscription) (*model.Subscription, error) {
	sub := input
	if err := sub.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetTopic(sub.TopicID); err != nil {
		return nil, model.NewValidationError("topic_id", "主题不存在")
	}
	if _, err := s.store.GetRecipient(sub.RecipientID); err != nil {
		return nil, model.NewValidationError("recipient_id", "接收人不存在")
	}
	now := time.Now()
	sub.ID = idgen.Hex()
	sub.CreatedAt = now
	sub.UpdatedAt = now
	if err := s.store.CreateSubscription(&sub); err != nil {
		return nil, err
	}
	_ = s.RecalculateTopicSubscribers(sub.TopicID)
	s.log.Infof("订阅创建成功: id=%s topic=%s recipient=%s", sub.ID, sub.TopicID, sub.RecipientID)
	return &sub, nil
}

// GetSubscription 获取订阅详情。
func (s *Service) GetSubscription(id string) (*model.Subscription, error) {
	return s.store.GetSubscription(id)
}

// ListSubscriptions 分页查询订阅。
func (s *Service) ListSubscriptions(filter model.SubscriptionFilter, page, size int) ([]*model.Subscription, int, error) {
	all := s.store.ListSubscriptions()
	matched := make([]*model.Subscription, 0, len(all))
	for _, sub := range all {
		if filter.Match(sub) {
			matched = append(matched, sub)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateSubscription 更新订阅渠道与状态。
func (s *Service) UpdateSubscription(id string, input model.Subscription) (*model.Subscription, error) {
	exist, err := s.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	if input.ChannelType != "" {
		exist.ChannelType = input.ChannelType
	}
	if input.Status != "" {
		exist.Status = input.Status
	}
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateSubscription(exist); err != nil {
		return nil, err
	}
	_ = s.RecalculateTopicSubscribers(exist.TopicID)
	return exist, nil
}

// DeleteSubscription 删除订阅。
func (s *Service) DeleteSubscription(id string) error {
	sub, err := s.store.GetSubscription(id)
	if err != nil {
		return err
	}
	if err := s.store.DeleteSubscription(id); err != nil {
		return err
	}
	_ = s.RecalculateTopicSubscribers(sub.TopicID)
	return nil
}

// Unsubscribe 退订：执行 subscribed→unsubscribed 状态机流转。
func (s *Service) Unsubscribe(id string) (*model.Subscription, error) {
	sub, err := s.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionSubscription(sub.Status, model.SubscriptionUnsubscribed) {
		return nil, model.NewValidationError("status", "订阅状态不允许从 "+sub.Status+" 流转到 "+model.SubscriptionUnsubscribed)
	}
	sub.Status = model.SubscriptionUnsubscribed
	sub.UpdatedAt = time.Now()
	if err := s.store.UpdateSubscription(sub); err != nil {
		return nil, err
	}
	_ = s.RecalculateTopicSubscribers(sub.TopicID)
	return sub, nil
}

// Subscribe 重新订阅。
func (s *Service) Subscribe(id string) (*model.Subscription, error) {
	sub, err := s.store.GetSubscription(id)
	if err != nil {
		return nil, err
	}
	sub.Status = model.SubscriptionSubscribed
	sub.UpdatedAt = time.Now()
	if err := s.store.UpdateSubscription(sub); err != nil {
		return nil, err
	}
	_ = s.RecalculateTopicSubscribers(sub.TopicID)
	return sub, nil
}

// BatchDeleteSubscriptions 批量删除订阅。
func (s *Service) BatchDeleteSubscriptions(ids []string) (int, error) {
	deleted := 0
	for _, id := range ids {
		if err := s.DeleteSubscription(id); err == nil {
			deleted++
		}
	}
	if deleted == 0 && len(ids) > 0 {
		return 0, store.ErrNotFound
	}
	return deleted, nil
}
