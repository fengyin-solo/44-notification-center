package store

import "notification/internal/model"

// CreateSendRecord 创建发送记录。
func (s *MemoryStore) CreateSendRecord(r *model.SendRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.sendRecords[r.ID] = r
	return nil
}

// GetSendRecord 按 ID 获取发送记录。
func (s *MemoryStore) GetSendRecord(id string) (*model.SendRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.sendRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return r, nil
}

// ListSendRecords 返回全部发送记录。
func (s *MemoryStore) ListSendRecords() []*model.SendRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SendRecord, 0, len(s.sendRecords))
	for _, r := range s.sendRecords {
		list = append(list, r)
	}
	return list
}

// UpdateSendRecord 更新发送记录。
func (s *MemoryStore) UpdateSendRecord(r *model.SendRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sendRecords[r.ID]; !ok {
		return ErrNotFound
	}
	s.sendRecords[r.ID] = r
	return nil
}

// DeleteSendRecord 删除发送记录。
func (s *MemoryStore) DeleteSendRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.sendRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.sendRecords, id)
	return nil
}
