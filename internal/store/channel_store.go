package store

import "notification/internal/model"

// CreateChannel 创建渠道，名称唯一。
func (s *MemoryStore) CreateChannel(c *model.Channel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.channels {
		if exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.channels[c.ID] = c
	return nil
}

// GetChannel 按 ID 获取渠道。
func (s *MemoryStore) GetChannel(id string) (*model.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.channels[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

// GetChannelByName 按名称获取渠道。
func (s *MemoryStore) GetChannelByName(name string) (*model.Channel, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.channels {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

// ListChannels 返回全部渠道。
func (s *MemoryStore) ListChannels() []*model.Channel {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Channel, 0, len(s.channels))
	for _, c := range s.channels {
		list = append(list, c)
	}
	return list
}

// UpdateChannel 更新渠道，名称唯一（排除自身）。
func (s *MemoryStore) UpdateChannel(c *model.Channel) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.channels[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.channels {
		if exist.ID != c.ID && exist.Name == c.Name {
			return ErrConflict
		}
	}
	s.channels[c.ID] = c
	return nil
}

// DeleteChannel 删除渠道。
func (s *MemoryStore) DeleteChannel(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.channels[id]; !ok {
		return ErrNotFound
	}
	delete(s.channels, id)
	return nil
}
