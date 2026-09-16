package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/idgen"
)

// CreateSchedule 创建定时任务，校验消息存在。
func (s *Service) CreateSchedule(input model.Schedule) (*model.Schedule, error) {
	sc := input
	if err := sc.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMessage(sc.MessageID); err != nil {
		return nil, model.NewValidationError("message_id", "消息不存在")
	}
	now := time.Now()
	sc.ID = idgen.Hex()
	sc.CreatedAt = now
	sc.UpdatedAt = now
	if err := s.store.CreateSchedule(&sc); err != nil {
		return nil, err
	}
	s.log.Infof("定时任务创建成功: id=%s message=%s", sc.ID, sc.MessageID)
	return &sc, nil
}

// GetSchedule 获取定时任务详情。
func (s *Service) GetSchedule(id string) (*model.Schedule, error) {
	return s.store.GetSchedule(id)
}

// ListSchedules 分页查询定时任务。
func (s *Service) ListSchedules(filter model.ScheduleFilter, page, size int) ([]*model.Schedule, int, error) {
	all := s.store.ListSchedules()
	matched := make([]*model.Schedule, 0, len(all))
	for _, sc := range all {
		if filter.Match(sc) {
			matched = append(matched, sc)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateSchedule 更新定时任务。
func (s *Service) UpdateSchedule(id string, input model.Schedule) (*model.Schedule, error) {
	exist, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	exist.CronExpr = input.CronExpr
	exist.NextRunAt = input.NextRunAt
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateSchedule(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteSchedule 删除定时任务。
func (s *Service) DeleteSchedule(id string) error {
	return s.store.DeleteSchedule(id)
}

// ExecuteSchedule 执行定时任务：pending→executed。
func (s *Service) ExecuteSchedule(id string) (*model.Schedule, error) {
	sc, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionSchedule(sc.Status, model.ScheduleExecuted) {
		return nil, model.NewValidationError("status", "定时任务状态不允许从 "+sc.Status+" 流转到 executed")
	}
	now := time.Now()
	sc.Status = model.ScheduleExecuted
	sc.LastRunAt = &now
	sc.UpdatedAt = now
	if err := s.store.UpdateSchedule(sc); err != nil {
		return nil, err
	}
	return sc, nil
}

// CancelSchedule 取消定时任务：pending→cancelled。
func (s *Service) CancelSchedule(id string) (*model.Schedule, error) {
	sc, err := s.store.GetSchedule(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionSchedule(sc.Status, model.ScheduleCancelled) {
		return nil, model.NewValidationError("status", "定时任务状态不允许从 "+sc.Status+" 流转到 cancelled")
	}
	sc.Status = model.ScheduleCancelled
	sc.UpdatedAt = time.Now()
	if err := s.store.UpdateSchedule(sc); err != nil {
		return nil, err
	}
	return sc, nil
}

// BatchDeleteSchedules 批量删除定时任务。
func (s *Service) BatchDeleteSchedules(ids []string) (int, error) {
	deleted := 0
	for _, id := range ids {
		if err := s.store.DeleteSchedule(id); err == nil {
			deleted++
		}
	}
	if deleted == 0 && len(ids) > 0 {
		return 0, store.ErrNotFound
	}
	return deleted, nil
}
