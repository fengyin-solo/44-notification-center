package store

import "notification/internal/model"

// CreateTemplate 创建模板，名称唯一。
func (s *MemoryStore) CreateTemplate(t *model.Template) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.templates {
		if exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.templates[t.ID] = t
	return nil
}

// GetTemplate 按 ID 获取模板。
func (s *MemoryStore) GetTemplate(id string) (*model.Template, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.templates[id]
	if !ok {
		return nil, ErrNotFound
	}
	return t, nil
}

// GetTemplateByName 按名称获取模板。
func (s *MemoryStore) GetTemplateByName(name string) (*model.Template, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, t := range s.templates {
		if t.Name == name {
			return t, nil
		}
	}
	return nil, ErrNotFound
}

// ListTemplates 返回全部模板。
func (s *MemoryStore) ListTemplates() []*model.Template {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Template, 0, len(s.templates))
	for _, t := range s.templates {
		list = append(list, t)
	}
	return list
}

// UpdateTemplate 更新模板，名称唯一（排除自身）。
func (s *MemoryStore) UpdateTemplate(t *model.Template) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.templates[t.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.templates {
		if exist.ID != t.ID && exist.Name == t.Name {
			return ErrConflict
		}
	}
	s.templates[t.ID] = t
	return nil
}

// DeleteTemplate 删除模板。
func (s *MemoryStore) DeleteTemplate(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.templates[id]; !ok {
		return ErrNotFound
	}
	delete(s.templates, id)
	return nil
}
