package service

import (
	"sort"

	"notification/internal/model"
)

// Overview 汇总全局统计概览。
func (s *Service) Overview() (*model.StatsOverview, error) {
	messages := s.store.ListMessages()
	records := s.store.ListSendRecords()
	topics := s.store.ListTopics()
	subscriptions := s.store.ListSubscriptions()

	o := &model.StatsOverview{
		Channels:        len(s.store.ListChannels()),
		Templates:       len(s.store.ListTemplates()),
		Topics:          len(topics),
		Recipients:      len(s.store.ListRecipients()),
		Subscriptions:   len(subscriptions),
		Messages:        len(messages),
		SendRecords:     len(records),
		RetryPolicies:   len(s.store.ListRetryPolicies()),
		Schedules:       len(s.store.ListSchedules()),
		TotalRecipients: len(s.store.ListRecipients()),
	}
	for _, m := range messages {
		switch m.Status {
		case model.MessageSent:
			o.MessagesSent++
		case model.MessageFailed:
			o.MessagesFailed++
		case model.MessagePending:
			o.MessagesPending++
		}
	}
	for _, r := range records {
		switch r.Status {
		case model.SendRecordSuccess:
			o.RecordsSucceeded++
		case model.SendRecordFailed:
			o.RecordsFailed++
		}
	}
	for _, sub := range subscriptions {
		if sub.Status == model.SubscriptionSubscribed {
			o.TotalSubscribers++
		}
	}
	total := o.RecordsSucceeded + o.RecordsFailed
	if total > 0 {
		o.SuccessRate = float64(o.RecordsSucceeded) / float64(total) * 100
	}
	return o, nil
}

// ChannelBreakdown 按渠道分组的发送统计。
func (s *Service) ChannelBreakdown() ([]*model.ChannelStats, error) {
	messages := s.store.ListMessages()
	records := s.store.ListSendRecords()
	byChannel := make(map[string]*model.ChannelStats)
	for _, m := range messages {
		cs := ensureChannelStat(byChannel, m.ChannelType)
		cs.Messages++
	}
	for _, r := range records {
		cs := ensureChannelStat(byChannel, r.ChannelType)
		cs.Records++
		switch r.Status {
		case model.SendRecordSuccess:
			cs.Succeeded++
		case model.SendRecordFailed:
			cs.Failed++
		}
	}
	out := make([]*model.ChannelStats, 0, len(byChannel))
	for _, cs := range byChannel {
		out = append(out, cs)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Messages > out[j].Messages
	})
	return out, nil
}

// MessageStatusBreakdown 消息按状态分组计数。
func (s *Service) MessageStatusBreakdown() ([]*model.MessageStatusStats, error) {
	messages := s.store.ListMessages()
	counts := make(map[string]int)
	for _, m := range messages {
		counts[m.Status]++
	}
	out := make([]*model.MessageStatusStats, 0, len(counts))
	for status, count := range counts {
		out = append(out, &model.MessageStatusStats{Status: status, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Status < out[j].Status
	})
	return out, nil
}

// MessagePriorityBreakdown 消息按优先级分组计数。
func (s *Service) MessagePriorityBreakdown() ([]*model.MessagePriorityStats, error) {
	messages := s.store.ListMessages()
	counts := make(map[string]int)
	for _, m := range messages {
		counts[m.Priority]++
	}
	out := make([]*model.MessagePriorityStats, 0, len(counts))
	for priority, count := range counts {
		out = append(out, &model.MessagePriorityStats{Priority: priority, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		return priorityRank(out[i].Priority) < priorityRank(out[j].Priority)
	})
	return out, nil
}

// TopTemplates 返回使用次数最高的前 N 个模板。
func (s *Service) TopTemplates(n int) ([]*model.TemplateUsage, error) {
	if n <= 0 {
		n = 5
	}
	messages := s.store.ListMessages()
	counts := make(map[string]int)
	for _, m := range messages {
		if m.TemplateID != "" {
			counts[m.TemplateID]++
		}
	}
	usages := make([]*model.TemplateUsage, 0, len(counts))
	for tid, count := range counts {
		name := ""
		if t, err := s.store.GetTemplate(tid); err == nil {
			name = t.Name
		}
		usages = append(usages, &model.TemplateUsage{TemplateID: tid, TemplateName: name, MessageCount: count})
	}
	sort.Slice(usages, func(i, j int) bool {
		return usages[i].MessageCount > usages[j].MessageCount
	})
	if len(usages) > n {
		usages = usages[:n]
	}
	return usages, nil
}

// TopicSubscriptionBreakdown 返回每个主题的订阅数统计。
func (s *Service) TopicSubscriptionBreakdown() ([]*model.TopicSubscriptionStats, error) {
	subscriptions := s.store.ListSubscriptions()
	stats := make(map[string]*model.TopicSubscriptionStats)
	for _, sub := range subscriptions {
		ts := ensureTopicStat(stats, sub.TopicID)
		ts.SubscriberCount++
		if sub.Status == model.SubscriptionSubscribed {
			ts.ActiveCount++
		}
	}
	for _, t := range s.store.ListTopics() {
		ts := ensureTopicStat(stats, t.ID)
		ts.TopicName = t.Name
	}
	out := make([]*model.TopicSubscriptionStats, 0, len(stats))
	for _, ts := range stats {
		out = append(out, ts)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].SubscriberCount > out[j].SubscriberCount
	})
	return out, nil
}

// ExportSnapshot 导出全部数据的汇总快照。
func (s *Service) ExportSnapshot() (map[string]interface{}, error) {
	overview, err := s.Overview()
	if err != nil {
		return nil, err
	}
	channels, err := s.ChannelBreakdown()
	if err != nil {
		return nil, err
	}
	statuses, err := s.MessageStatusBreakdown()
	if err != nil {
		return nil, err
	}
	priorities, err := s.MessagePriorityBreakdown()
	if err != nil {
		return nil, err
	}
	topTemplates, err := s.TopTemplates(10)
	if err != nil {
		return nil, err
	}
	topicStats, err := s.TopicSubscriptionBreakdown()
	if err != nil {
		return nil, err
	}
	snapshot := map[string]interface{}{
		"overview":       overview,
		"channels":       s.store.ListChannels(),
		"templates":      s.store.ListTemplates(),
		"topics":         s.store.ListTopics(),
		"recipients":     s.store.ListRecipients(),
		"subscriptions":  s.store.ListSubscriptions(),
		"messages":       s.store.ListMessages(),
		"send_records":   s.store.ListSendRecords(),
		"retry_policies": s.store.ListRetryPolicies(),
		"schedules":      s.store.ListSchedules(),
		"stats": map[string]interface{}{
			"channel_breakdown":   channels,
			"message_status":      statuses,
			"message_priority":    priorities,
			"top_templates":       topTemplates,
			"topic_subscriptions": topicStats,
		},
	}
	return snapshot, nil
}

// ensureChannelStat 获取或创建渠道统计对象。
func ensureChannelStat(m map[string]*model.ChannelStats, channelType string) *model.ChannelStats {
	cs, ok := m[channelType]
	if !ok {
		cs = &model.ChannelStats{ChannelType: channelType}
		m[channelType] = cs
	}
	return cs
}

// ensureTopicStat 获取或创建主题统计对象。
func ensureTopicStat(m map[string]*model.TopicSubscriptionStats, topicID string) *model.TopicSubscriptionStats {
	ts, ok := m[topicID]
	if !ok {
		ts = &model.TopicSubscriptionStats{TopicID: topicID}
		m[topicID] = ts
	}
	return ts
}
