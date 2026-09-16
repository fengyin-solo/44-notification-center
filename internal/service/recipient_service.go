package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/idgen"
)

// CreateRecipient 创建接收人。
func (s *Service) CreateRecipient(input model.Recipient) (*model.Recipient, error) {
	r := input
	if err := r.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	r.ID = idgen.Hex()
	r.CreatedAt = now
	r.UpdatedAt = now
	if err := s.store.CreateRecipient(&r); err != nil {
		return nil, err
	}
	s.log.Infof("接收人创建成功: id=%s name=%s", r.ID, r.Name)
	return &r, nil
}

// GetRecipient 获取接收人详情。
func (s *Service) GetRecipient(id string) (*model.Recipient, error) {
	return s.store.GetRecipient(id)
}

// ListRecipients 分页查询接收人。
func (s *Service) ListRecipients(filter model.RecipientFilter, page, size int) ([]*model.Recipient, int, error) {
	all := s.store.ListRecipients()
	matched := make([]*model.Recipient, 0, len(all))
	for _, r := range all {
		if filter.Match(r) {
			matched = append(matched, r)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateRecipient 更新接收人。
func (s *Service) UpdateRecipient(id string, input model.Recipient) (*model.Recipient, error) {
	exist, err := s.store.GetRecipient(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.ChannelType = input.ChannelType
	exist.Address = input.Address
	exist.Group = input.Group
	if input.Status != "" {
		exist.Status = input.Status
	}
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateRecipient(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteRecipient 删除接收人。
func (s *Service) DeleteRecipient(id string) error {
	return s.store.DeleteRecipient(id)
}

// SetRecipientStatus 更新接收人状态（active/unsubscribed）。
func (s *Service) SetRecipientStatus(id, status string) (*model.Recipient, error) {
	r, err := s.store.GetRecipient(id)
	if err != nil {
		return nil, err
	}
	if status != model.RecipientActive && status != model.RecipientUnsubscribed {
		return nil, model.NewValidationError("status", "接收人状态不合法")
	}
	r.Status = status
	r.UpdatedAt = time.Now()
	if err := s.store.UpdateRecipient(r); err != nil {
		return nil, err
	}
	return r, nil
}

// BatchDeleteRecipients 批量删除接收人。
func (s *Service) BatchDeleteRecipients(ids []string) (int, error) {
	deleted := 0
	for _, id := range ids {
		if err := s.store.DeleteRecipient(id); err == nil {
			deleted++
		}
	}
	if deleted == 0 && len(ids) > 0 {
		return 0, store.ErrNotFound
	}
	return deleted, nil
}
