package model

// CountGroup 按某个维度分组的计数结果。
type CountGroup struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

// StatsOverview 全局统计概览。
type StatsOverview struct {
	Channels         int     `json:"channels"`
	Templates        int     `json:"templates"`
	Topics           int     `json:"topics"`
	Recipients       int     `json:"recipients"`
	Subscriptions    int     `json:"subscriptions"`
	Messages         int     `json:"messages"`
	SendRecords      int     `json:"send_records"`
	RetryPolicies    int     `json:"retry_policies"`
	Schedules        int     `json:"schedules"`
	MessagesSent     int     `json:"messages_sent"`
	MessagesFailed   int     `json:"messages_failed"`
	MessagesPending  int     `json:"messages_pending"`
	RecordsSucceeded int     `json:"records_succeeded"`
	RecordsFailed    int     `json:"records_failed"`
	TotalSubscribers int     `json:"total_subscribers"`
	TotalRecipients  int     `json:"total_recipients"`
	SuccessRate      float64 `json:"success_rate"`
}

// ChannelStats 按渠道分组统计。
type ChannelStats struct {
	ChannelType string `json:"channel_type"`
	Messages    int    `json:"messages"`
	Records     int    `json:"records"`
	Succeeded   int    `json:"succeeded"`
	Failed      int    `json:"failed"`
}

// MessageStatusStats 消息按状态分组统计。
type MessageStatusStats struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

// MessagePriorityStats 消息按优先级分组统计。
type MessagePriorityStats struct {
	Priority string `json:"priority"`
	Count    int    `json:"count"`
}

// TemplateUsage 模板使用次数统计（TOP N）。
type TemplateUsage struct {
	TemplateID   string `json:"template_id"`
	TemplateName string `json:"template_name"`
	MessageCount int    `json:"message_count"`
}

// TopicSubscriptionStats 主题订阅数统计。
type TopicSubscriptionStats struct {
	TopicID         string `json:"topic_id"`
	TopicName       string `json:"topic_name"`
	SubscriberCount int    `json:"subscriber_count"`
	ActiveCount     int    `json:"active_count"`
}
