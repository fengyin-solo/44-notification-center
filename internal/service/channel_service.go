package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/idgen"
)

// CreateChannel 创建通知渠道。
func (s *Service) CreateChannel(input model.Channel) (*model.Channel, error) {
	c := input
	if err := c.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	c.ID = idgen.Hex()
	c.CreatedAt = now
	c.UpdatedAt = now
	if err := s.store.CreateChannel(&c); err != nil {
		return nil, err
	}
	s.log.Infof("渠道创建成功: id=%s name=%s", c.ID, c.Name)
	return &c, nil
}

// GetChannel 获取渠道详情。
func (s *Service) GetChannel(id string) (*model.Channel, error) {
	return s.store.GetChannel(id)
}

// ListChannels 分页查询渠道。
func (s *Service) ListChannels(filter model.ChannelFilter, page, size int) ([]*model.Channel, int, error) {
	all := s.store.ListChannels()
	matched := make([]*model.Channel, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		if matched[i].Priority != matched[j].Priority {
			return matched[i].Priority > matched[j].Priority
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateChannel 更新渠道。
func (s *Service) UpdateChannel(id string, input model.Channel) (*model.Channel, error) {
	exist, err := s.store.GetChannel(id)
	if err != nil {
		return nil, err
	}
	exist.Name = input.Name
	exist.Type = input.Type
	exist.Config = input.Config
	exist.Priority = input.Priority
	if input.Status != "" {
		exist.Status = input.Status
	}
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateChannel(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteChannel 删除渠道。
func (s *Service) DeleteChannel(id string) error {
	if err := s.store.DeleteChannel(id); err != nil {
		return err
	}
	s.log.Infof("渠道删除成功: id=%s", id)
	return nil
}

// BatchDeleteChannels 批量删除渠道。
func (s *Service) BatchDeleteChannels(ids []string) (int, error) {
	deleted := 0
	for _, id := range ids {
		if err := s.store.DeleteChannel(id); err == nil {
			deleted++
		}
	}
	if deleted == 0 && len(ids) > 0 {
		return 0, store.ErrNotFound
	}
	return deleted, nil
}
