package store

import "notification/internal/model"

// CreateTopic 创建主题，名称唯一。
func (s *MemoryStore) CreateTopic(t *model.Topic) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.topics {
		if exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.topics[t.ID] = t
	return nil
}

// GetTopic 按 ID 获取主题。
func (s *MemoryStore) GetTopic(id string) (*model.Topic, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.topics[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

// ListTopics 返回全部主题。
func (s *MemoryStore) ListTopics() []*model.Topic {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Topic, 0, len(s.topics))
	for _, t := range s.topics {
		list = append(list, t)
	}
	return list
}

// UpdateTopic 更新主题，名称唯一（排除自身）。
func (s *MemoryStore) UpdateTopic(t *model.Topic) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.topics[t.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.topics {
		if exist.ID != t.ID && exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.topics[t.ID] = t
	return nil
}

// DeleteTopic 删除主题。
func (s *MemoryStore) DeleteTopic(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.topics[id]; !ok {
		return ErrNotFound
	}
	delete(s.topics, id)
	return nil
}
