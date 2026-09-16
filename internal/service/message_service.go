package service

import (
	"sort"
	"time"

	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/idgen"
)

// CreateMessage 创建消息，校验模板与主题存在。
func (s *Service) CreateMessage(input model.Message) (*model.Message, error) {
	m := input
	if err := m.Validate(); err != nil {
		return nil, err
	}
	if m.TemplateID != "" {
		if _, err := s.store.GetTemplate(m.TemplateID); err != nil {
			return nil, model.NewValidationError("template_id", "模板不存在")
		}
	}
	if m.TopicID != "" {
		if _, err := s.store.GetTopic(m.TopicID); err != nil {
			return nil, model.NewValidationError("topic_id", "主题不存在")
		}
	}
	now := time.Now()
	m.ID = idgen.Hex()
	m.CreatedAt = now
	m.UpdatedAt = now
	if err := s.store.CreateMessage(&m); err != nil {
		return nil, err
	}
	s.log.Infof("消息创建成功: id=%s title=%s", m.ID, m.Title)
	return &m, nil
}

// GetMessage 获取消息详情。
func (s *Service) GetMessage(id string) (*model.Message, error) {
	return s.store.GetMessage(id)
}

// ListMessages 分页查询消息。
func (s *Service) ListMessages(filter model.MessageFilter, page, size int) ([]*model.Message, int, error) {
	all := s.store.ListMessages()
	matched := make([]*model.Message, 0, len(all))
	for _, m := range all {
		if filter.Match(m) {
			matched = append(matched, m)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		pi := priorityRank(matched[i].Priority)
		pj := priorityRank(matched[j].Priority)
		if pi != pj {
			return pi < pj
		}
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	return paginate(matched, page, size), total, nil
}

// UpdateMessage 更新草稿消息。
func (s *Service) UpdateMessage(id string, input model.Message) (*model.Message, error) {
	exist, err := s.store.GetMessage(id)
	if err != nil {
		return nil, err
	}
	if exist.Status != model.MessageDraft && exist.Status != model.MessagePending {
		return nil, model.NewValidationError("status", "仅草稿或待发送状态的消息可编辑")
	}
	exist.Title = input.Title
	exist.Content = input.Content
	exist.ChannelType = input.ChannelType
	exist.Priority = input.Priority
	if err := exist.Validate(); err != nil {
		return nil, err
	}
	exist.UpdatedAt = time.Now()
	if err := s.store.UpdateMessage(exist); err != nil {
		return nil, err
	}
	return exist, nil
}

// DeleteMessage 删除消息。
func (s *Service) DeleteMessage(id string) error {
	return s.store.DeleteMessage(id)
}

// TransitionMessage 执行消息状态机流转。
func (s *Service) TransitionMessage(id, to string) (*model.Message, error) {
	m, err := s.store.GetMessage(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionMessage(m.Status, to) {
		return nil, model.NewValidationError("status", "消息状态不允许从 "+m.Status+" 流转到 "+to)
	}
	m.Status = to
	m.UpdatedAt = time.Now()
	if to == model.MessageSent {
		now := time.Now()
		m.SentAt = &now
	}
	if err := s.store.UpdateMessage(m); err != nil {
		return nil, err
	}
	return m, nil
}

// SendMessage 发送消息：流转状态并生成发送记录。
func (s *Service) SendMessage(id string) (*model.Message, int, error) {
	m, err := s.store.GetMessage(id)
	if err != nil {
		return nil, 0, err
	}
	if m.Status == model.MessageDraft {
		m.Status = model.MessagePending
		m.UpdatedAt = time.Now()
		if err := s.store.UpdateMessage(m); err != nil {
			return nil, 0, err
		}
	}
	if m.Status == model.MessagePending {
		m.Status = model.MessageSending
		m.UpdatedAt = time.Now()
		if err := s.store.UpdateMessage(m); err != nil {
			return nil, 0, err
		}
	}
	if m.Status != model.MessageSending {
		return nil, 0, model.NewValidationError("status", "消息当前状态为 "+m.Status+"，无法发送")
	}
	recipients := s.selectRecipients(m)
	recordCount := 0
	now := time.Now()
	for _, r := range recipients {
		rec := &model.SendRecord{
			ID:          idgen.Hex(),
			MessageID:   m.ID,
			RecipientID: r.ID,
			ChannelType: m.ChannelType,
			Status:      model.SendRecordSuccess,
			Attempts:    1,
			DurationMs:  int64(10 + time.Now().Nanosecond()%50),
			SentAt:      &now,
			CreatedAt:   now,
		}
		if err := s.store.CreateSendRecord(rec); err == nil {
			recordCount++
		}
	}
	m.Status = model.MessageSent
	m.UpdatedAt = time.Now()
	m.SentAt = &now
	if err := s.store.UpdateMessage(m); err != nil {
		return nil, recordCount, err
	}
	s.log.Infof("消息发送完成: id=%s records=%d", m.ID, recordCount)
	return m, recordCount, nil
}

// BatchSendMessages 批量发送消息，返回成功发送的数量。
func (s *Service) BatchSendMessages(ids []string) (int, error) {
	sent := 0
	for _, id := range ids {
		if _, _, err := s.SendMessage(id); err == nil {
			sent++
		}
	}
	if sent == 0 && len(ids) > 0 {
		return 0, model.NewValidationError("ids", "没有消息被成功发送")
	}
	return sent, nil
}

// BatchUpdateMessageStatus 批量更新消息状态（执行状态机校验）。
func (s *Service) BatchUpdateMessageStatus(ids []string, to string) (int, error) {
	updated := 0
	for _, id := range ids {
		if _, err := s.TransitionMessage(id, to); err == nil {
			updated++
		}
	}
	if updated == 0 && len(ids) > 0 {
		return 0, model.NewValidationError("status", "没有消息完成状态更新")
	}
	return updated, nil
}

// BatchDeleteMessages 批量删除消息。
func (s *Service) BatchDeleteMessages(ids []string) (int, error) {
	deleted := 0
	for _, id := range ids {
		if err := s.store.DeleteMessage(id); err == nil {
			deleted++
		}
	}
	if deleted == 0 && len(ids) > 0 {
		return 0, store.ErrNotFound
	}
	return deleted, nil
}

// selectRecipients 选择消息的目标接收人（优先订阅，退化为渠道匹配）。
func (s *Service) selectRecipients(m *model.Message) []*model.Recipient {
	chosen := make(map[string]bool)
	result := make([]*model.Recipient, 0)
	for _, sub := range s.store.ListSubscriptions() {
		if sub.TopicID == m.TopicID && sub.Status == model.SubscriptionSubscribed {
			if r, err := s.store.GetRecipient(sub.RecipientID); err == nil && r.Status == model.RecipientActive {
				if !chosen[r.ID] {
					chosen[r.ID] = true
					result = append(result, r)
				}
			}
		}
	}
	if len(result) == 0 {
		for _, r := range s.store.ListRecipients() {
			if r.ChannelType == m.ChannelType && r.Status == model.RecipientActive {
				if !chosen[r.ID] {
					chosen[r.ID] = true
					result = append(result, r)
				}
			}
		}
	}
	return result
}

// priorityRank 将优先级映射为排序权重（数值越小越优先）。
func priorityRank(p string) int {
	switch p {
	case model.PriorityUrgent:
		return 0
	case model.PriorityHigh:
		return 1
	case model.PriorityNormal:
		return 2
	case model.PriorityLow:
		return 3
	default:
		return 4
	}
}
