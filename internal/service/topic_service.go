package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/idgen"
)

// CreateTopic 创建订阅主题。
func (s *Service) CreateTopic(input model.Topic) (*model.Topic, error) {
	t := input
	if err := t.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	t.ID = idgen.Hex()
	t.SubscriberCount = 0
	t.CreatedAt = now
	t.UpdatedAt = now
	if err := s.store.CreateTopic(&t); err != nil {
		return nil, err
	}
	s.log.Infof("主题创建成功: id=%s name=%s", t.ID, t.Name)
	return &t, nil
}

// GetTopic 获取主题详情。
func (s *Service) GetTopic(id string) (*model.Topic, error) {
	return s.store.GetTopic(id)
}

// ListTopics 分页查询主题。
func (s *Service) ListTopics(filter model.TopicFilter, page, size int) ([]*model.Topic, int, error) {
	all := s.store.ListTopics()
	matched := make([]*model.Topic, 0, len(all))
	for _, t := range all {
		if filter.Match(t) {
			matched = append(matched, t)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateTopic 更新主题。
func (s *Service) UpdateTopic(id string, input model.Topic) (*model.Topic, error) {
	exist, err := s.store.GetTopic(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Description = input.Description
	if input.Status != "" {
		exist.Status = input.Status
	}
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateTopic(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteTopic 删除主题。
func (s *Service) DeleteTopic(id string) error {
	return s.store.DeleteTopic(id)
}

// RecalculateTopicSubscribers 重新计算主题订阅数并持久化。
func (s *Service) RecalculateTopicSubscribers(topicID string) error {
	t, err := s.store.GetTopic(topicID)
	if err != nil {
		return err
	}
	count := 0
	for _, sub := range s.store.ListSubscriptions() {
		if sub.TopicID == topicID && sub.Status == model.SubscriptionSubscribed {
			count++
		}
	}
	t.SubscriberCount = count
	t.UpdatedAt = time.Now()
	return s.store.UpdateTopic(t)
}

// BatchDeleteTopics 批量删除主题。
func (s *Service) BatchDeleteTopics(ids []string) (int, error) {
	deleted := 0
	for _, id := range ids {
		if err := s.store.DeleteTopic(id); err == nil {
			deleted++
		}
	}
	if deleted == 0 && len(ids) > 0 {
		return 0, store.ErrNotFound
	}
	return deleted, nil
}
