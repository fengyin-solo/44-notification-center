package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/idgen"
)

// CreateRetryPolicy 创建重试策略。
func (s *Service) CreateRetryPolicy(input model.RetryPolicy) (*model.RetryPolicy, error) {
	p := input
	if err := p.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	p.ID = idgen.Hex()
	p.CreatedAt = now
	p.UpdatedAt = now
	if err := s.store.CreateRetryPolicy(&p); err != nil {
		return nil, err
	}
	s.log.Infof("重试策略创建成功: id=%s name=%s", p.ID, p.Name)
	return &p, nil
}

// GetRetryPolicy 获取重试策略详情。
func (s *Service) GetRetryPolicy(id string) (*model.RetryPolicy, error) {
	return s.store.GetRetryPolicy(id)
}

// ListRetryPolicies 分页查询重试策略。
func (s *Service) ListRetryPolicies(filter model.RetryPolicyFilter, page, size int) ([]*model.RetryPolicy, int, error) {
	all := s.store.ListRetryPolicies()
	matched := make([]*model.RetryPolicy, 0, len(all))
	for _, p := range all {
		if filter.Match(p) {
			matched = append(matched, p)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateRetryPolicy 更新重试策略。
func (s *Service) UpdateRetryPolicy(id string, input model.RetryPolicy) (*model.RetryPolicy, error) {
	exist, err := s.store.GetRetryPolicy(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.ChannelType = input.ChannelType
	exist.MaxAttempts = input.MaxAttempts
	exist.BackoffMs = input.BackoffMs
	if input.Status != "" {
		exist.Status = input.Status
	}
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateRetryPolicy(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteRetryPolicy 删除重试策略。
func (s *Service) DeleteRetryPolicy(id string) error {
	return s.store.DeleteRetryPolicy(id)
}

// BatchDeleteRetryPolicies 批量删除重试策略。
func (s *Service) BatchDeleteRetryPolicies(ids []string) (int, error) {
	deleted := 0
	for _, id := range ids {
		if err := s.store.DeleteRetryPolicy(id); err == nil {
			deleted++
		}
	}
	if deleted == 0 && len(ids) > 0 {
		return 0, store.ErrNotFound
	}
	return deleted, nil
}
