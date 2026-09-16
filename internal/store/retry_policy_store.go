package store

import "notification/internal/model"

// CreateRetryPolicy 创建重试策略，名称唯一。
func (s *MemoryStore) CreateRetryPolicy(p *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.retryPolicies {
		if exist.Name == p.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[p.ID] = p
	return nil
}

// GetRetryPolicy 按 ID 获取重试策略。
func (s *MemoryStore) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.retryPolicies[id]
	if !ok {
		return nil, ErrNotFound
	}
	return p, nil
}

// GetRetryPolicyByName 按名称获取重试策略。
func (s *MemoryStore) GetRetryPolicyByName(name string) (*model.RetryPolicy, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, p := range s.retryPolicies {
		if p.Name == name {
			return p, nil
		}
	}
	return nil, ErrNotFound
}

// ListRetryPolicies 返回全部重试策略。
func (s *MemoryStore) ListRetryPolicies() []*model.RetryPolicy {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.RetryPolicy, 0, len(s.retryPolicies))
	for _, p := range s.retryPolicies {
		list = append(list, p)
	}
	return list
}

// UpdateRetryPolicy 更新重试策略，名称唯一（排除自身）。
func (s *MemoryStore) UpdateRetryPolicy(p *model.RetryPolicy) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[p.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.retryPolicies {
		if exist.ID != p.ID && exist.Name == p.Name {
			return ErrConflict
		}
	}
	s.retryPolicies[p.ID] = p
	return nil
}

// DeleteRetryPolicy 删除重试策略。
func (s *MemoryStore) DeleteRetryPolicy(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.retryPolicies[id]; !ok {
		return ErrNotFound
	}
	delete(s.retryPolicies, id)
	return nil
}
