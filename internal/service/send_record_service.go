package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/pkg/idgen"
)

// CreateSendRecord 创建发送记录，校验消息与接收人均存在。
func (s *Service) CreateSendRecord(input model.SendRecord) (*model.SendRecord, error) {
	r := input
	if err := r.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetMessage(r.MessageID); err != nil {
		return nil, model.NewValidationError("message_id", "消息不存在")
	}
	if _, err := s.store.GetRecipient(r.RecipientID); err != nil {
		return nil, model.NewValidationError("recipient_id", "接收人不存在")
	}
	now := time.Now()
	r.ID = idgen.Hex()
	r.CreatedAt = now
	if err := s.store.CreateSendRecord(&r); err != nil {
		return nil, err
	}
	return &r, nil
}

// GetSendRecord 获取发送记录详情。
func (s *Service) GetSendRecord(id string) (*model.SendRecord, error) {
	return s.store.GetSendRecord(id)
}

// ListSendRecords 分页查询发送记录。
func (s *Service) ListSendRecords(filter model.SendRecordFilter, page, size int) ([]*model.SendRecord, int, error) {
	all := s.store.ListSendRecords()
	matched := make([]*model.SendRecord, 0, len(all))
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

// UpdateSendRecord 更新发送记录（用于补录结果）。
func (s *Service) UpdateSendRecord(id string, input model.SendRecord) (*model.SendRecord, error) {
	exist, err := s.store.GetSendRecord(id)
	if err != nil {
		return nil, err
	}
	exist.Status = input.Status
	exist.Error = input.Error
	exist.DurationMs = input.DurationMs
	exist.Attempts = input.Attempts
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	if input.SentAt != nil {
		exist.SentAt = input.SentAt
	}
	if err := s.store.UpdateSendRecord(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteSendRecord 删除发送记录。
func (s *Service) DeleteSendRecord(id string) error {
	return s.store.DeleteSendRecord(id)
}

// MarkRecordSuccess 将记录标记为成功（pending/retrying→success）。
func (s *Service) MarkRecordSuccess(id string, durationMs int64) (*model.SendRecord, error) {
	r, err := s.store.GetSendRecord(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionSendRecord(r.Status, model.SendRecordSuccess) {
		return nil, model.NewValidationError("status", "发送记录状态不允许从 "+r.Status+" 流转到 success")
	}
	now := time.Now()
	r.Status = model.SendRecordSuccess
	r.DurationMs = durationMs
	r.SentAt = &now
	if err := s.store.UpdateSendRecord(r); err != nil {
		return nil, err
	}
	return r, nil
}

// MarkRecordFailed 将记录标记为失败（pending/retrying→failed）。
func (s *Service) MarkRecordFailed(id string, reason string) (*model.SendRecord, error) {
	r, err := s.store.GetSendRecord(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionSendRecord(r.Status, model.SendRecordFailed) {
		return nil, model.NewValidationError("status", "发送记录状态不允许从 "+r.Status+" 流转到 failed")
	}
	r.Status = model.SendRecordFailed
	r.Error = reason
	if err := s.store.UpdateSendRecord(r); err != nil {
		return nil, err
	}
	return r, nil
}

// RetryRecord 将失败记录转入重试（failed→retrying）。
func (s *Service) RetryRecord(id string) (*model.SendRecord, error) {
	r, err := s.store.GetSendRecord(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionSendRecord(r.Status, model.SendRecordRetrying) {
		return nil, model.NewValidationError("status", "发送记录状态不允许从 "+r.Status+" 流转到 retrying")
	}
	r.Status = model.SendRecordRetrying
	r.Attempts++
	if err := s.store.UpdateSendRecord(r); err != nil {
		return nil, err
	}
	return r, nil
}

// BatchUpdateRecordStatus 批量更新发送记录状态。
func (s *Service) BatchUpdateRecordStatus(ids []string, to string, reason string) (int, error) {
	updated := 0
	for _, id := range ids {
		var err error
		switch to {
		case model.SendRecordSuccess:
			_, err = s.MarkRecordSuccess(id, 0)
		case model.SendRecordFailed:
			_, err = s.MarkRecordFailed(id, reason)
		case model.SendRecordRetrying:
			_, err = s.RetryRecord(id)
		default:
			err = model.NewValidationError("status", "不支持的批量状态 "+to)
		}
		if err == nil {
			updated++
		}
	}
	return updated, nil
}
