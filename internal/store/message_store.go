package store

import "notification/internal/model"

// CreateMessage 创建消息。
func (s *MemoryStore) CreateMessage(m *model.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messages[m.ID] = m
	return nil
}

// GetMessage 按 ID 获取消息。
func (s *MemoryStore) GetMessage(id string) (*model.Message, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m, ok := s.messages[id]
	if !ok {
		return nil, ErrNotFound
	}
	return m, nil
}

// ListMessages 返回全部消息。
func (s *MemoryStore) ListMessages() []*model.Message {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Message, 0, len(s.messages))
	for _, m := range s.messages {
		list = append(list, m)
	}
	return list
}

// UpdateMessage 更新消息。
func (s *MemoryStore) UpdateMessage(m *model.Message) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messages[m.ID]; !ok {
		return ErrNotFound
	}
	s.messages[m.ID] = m
	return nil
}

// DeleteMessage 删除消息。
func (s *MemoryStore) DeleteMessage(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.messages[id]; !ok {
		return ErrNotFound
	}
	delete(s.messages, id)
	return nil
}
