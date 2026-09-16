package store

import "notification/internal/model"

// CreateRecipient 创建接收人，地址+渠道类型唯一。
func (s *MemoryStore) CreateRecipient(r *model.Recipient) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.recipients {
		if exist.ChannelType == r.ChannelType && exist.Address == r.Address {
			return ErrConflict
		}
	}
	s.recipients[r.ID] = r
	return nil
}

// GetRecipient 按 ID 获取接收人。
func (s *MemoryStore) GetRecipient(id string) (*model.Recipient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.recipients[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

// GetRecipientByAddress 按渠道类型+地址获取接收人。
func (s *MemoryStore) GetRecipientByAddress(channelType, address string) (*model.Recipient, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, r := range s.recipients {
		if r.ChannelType == channelType && r.Address == address {
			return r, nil
		}
	}
	return nil, ErrNotFound
}

// ListRecipients 返回全部接收人。
func (s *MemoryStore) ListRecipients() []*model.Recipient {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Recipient, 0, len(s.recipients))
	for _, r := range s.recipients {
		list = append(list, r)
	}
	return list
}

// UpdateRecipient 更新接收人，地址+渠道类型唯一（排除自身）。
func (s *MemoryStore) UpdateRecipient(r *model.Recipient) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recipients[r.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.recipients {
		if exist.ID != r.ID && exist.ChannelType == r.ChannelType && exist.Address == r.Address {
			return ErrConflict
		}
	}
	s.recipients[r.ID] = r
	return nil
}

// DeleteRecipient 删除接收人。
func (s *MemoryStore) DeleteRecipient(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.recipients[id]; !ok {
		return ErrNotFound
	}
	delete(s.recipients, id)
	return nil
}
