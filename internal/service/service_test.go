package service

import (
	"testing"

	"notification/internal/config"
	"notification/internal/model"
	"notification/internal/store"
	"notification/pkg/logger"
)

func newTestService() *Service {
	cfg := &config.Config{MaxPageSize: 100}
	log := logger.NewLevel(logger.LevelError)
	return New(store.NewMemoryStore(), log, cfg)
}

func TestMessageStateMachine(t *testing.T) {
	svc := newTestService()
	m, err := svc.CreateMessage(model.Message{Title: "测试", Content: "内容", ChannelType: model.ChannelTypeEmail})
	if err != nil {
		t.Fatalf("创建消息失败: %v", err)
	}
	if m.Status != model.MessageDraft {
		t.Fatalf("期望初始状态 draft，得到 %s", m.Status)
	}
	// draft -> pending -> sending -> sent
	for _, target := range []string{model.MessagePending, model.MessageSending, model.MessageSent} {
		m, err = svc.TransitionMessage(m.ID, target)
		if err != nil {
			t.Fatalf("流转到 %s 失败: %v", target, err)
		}
	}
	if m.Status != model.MessageSent {
		t.Fatalf("期望状态 sent，得到 %s", m.Status)
	}
	// sent 不能再流转
	if _, err := svc.TransitionMessage(m.ID, model.MessageFailed); err == nil {
		t.Fatalf("期望 sent 状态不可流转到 failed")
	}
}

func TestMessageInvalidTransition(t *testing.T) {
	svc := newTestService()
	m, _ := svc.CreateMessage(model.Message{Title: "测试", Content: "内容", ChannelType: model.ChannelTypeEmail})
	// draft 不能直接到 sent
	if _, err := svc.TransitionMessage(m.ID, model.MessageSent); err == nil {
		t.Fatalf("期望 draft->sent 非法流转报错")
	}
}

func TestTemplateStateMachine(t *testing.T) {
	svc := newTestService()
	tm, err := svc.CreateTemplate(model.Template{Name: "模板", Type: model.ChannelTypeEmail, Subject: "s", Content: "c"})
	if err != nil {
		t.Fatalf("创建模板失败: %v", err)
	}
	if _, err := svc.TransitionTemplate(tm.ID, model.TemplateActive); err != nil {
		t.Fatalf("draft->active 失败: %v", err)
	}
	if _, err := svc.TransitionTemplate(tm.ID, model.TemplateArchived); err != nil {
		t.Fatalf("active->archived 失败: %v", err)
	}
	// archived 不能再流转
	if _, err := svc.TransitionTemplate(tm.ID, model.TemplateActive); err == nil {
		t.Fatalf("期望 archived 状态不可流转")
	}
}

func TestSubscriptionForeignKeyValidation(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateSubscription(model.Subscription{TopicID: "nope", RecipientID: "nope", ChannelType: model.ChannelTypeEmail}); err == nil {
		t.Fatalf("期望主题不存在时报错")
	}
	// 创建主题与接收人后应成功
	tp, _ := svc.CreateTopic(model.Topic{Name: "主题"})
	rp, _ := svc.CreateRecipient(model.Recipient{Name: "张三", ChannelType: model.ChannelTypeEmail, Address: "a@b.com"})
	sub, err := svc.CreateSubscription(model.Subscription{TopicID: tp.ID, RecipientID: rp.ID, ChannelType: model.ChannelTypeEmail})
	if err != nil {
		t.Fatalf("创建订阅失败: %v", err)
	}
	// 主题订阅数应更新为 1
	got, _ := svc.GetTopic(tp.ID)
	if got.SubscriberCount != 1 {
		t.Fatalf("期望订阅数 1，得到 %d", got.SubscriberCount)
	}
	// 退订后订阅数归零
	if _, err := svc.Unsubscribe(sub.ID); err != nil {
		t.Fatalf("退订失败: %v", err)
	}
	got, _ = svc.GetTopic(tp.ID)
	if got.SubscriberCount != 0 {
		t.Fatalf("期望订阅数 0，得到 %d", got.SubscriberCount)
	}
}

func TestMessageForeignKeyValidation(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateMessage(model.Message{TemplateID: "missing", Title: "t", Content: "c", ChannelType: model.ChannelTypeEmail}); err == nil {
		t.Fatalf("期望模板不存在时报错")
	}
	if _, err := svc.CreateMessage(model.Message{TopicID: "missing", Title: "t", Content: "c", ChannelType: model.ChannelTypeEmail}); err == nil {
		t.Fatalf("期望主题不存在时报错")
	}
}

func TestSendMessageGeneratesRecords(t *testing.T) {
	svc := newTestService()
	tp, _ := svc.CreateTopic(model.Topic{Name: "主题"})
	rp, _ := svc.CreateRecipient(model.Recipient{Name: "张三", ChannelType: model.ChannelTypeEmail, Address: "a@b.com"})
	if _, err := svc.CreateSubscription(model.Subscription{TopicID: tp.ID, RecipientID: rp.ID, ChannelType: model.ChannelTypeEmail}); err != nil {
		t.Fatalf("创建订阅失败: %v", err)
	}
	m, _ := svc.CreateMessage(model.Message{TopicID: tp.ID, Title: "标题", Content: "内容", ChannelType: model.ChannelTypeEmail})
	sent, count, err := svc.SendMessage(m.ID)
	if err != nil {
		t.Fatalf("发送失败: %v", err)
	}
	if sent.Status != model.MessageSent {
		t.Fatalf("期望消息 sent，得到 %s", sent.Status)
	}
	if count != 1 {
		t.Fatalf("期望生成 1 条发送记录，得到 %d", count)
	}
}

func TestOverviewStats(t *testing.T) {
	svc := newTestService()
	svc.CreateChannel(model.Channel{Name: "邮件", Type: model.ChannelTypeEmail})
	svc.CreateTemplate(model.Template{Name: "模板", Type: model.ChannelTypeEmail, Subject: "s", Content: "c"})
	svc.CreateTopic(model.Topic{Name: "主题"})
	svc.CreateRecipient(model.Recipient{Name: "张三", ChannelType: model.ChannelTypeEmail, Address: "a@b.com"})
	svc.CreateMessage(model.Message{Title: "标题", Content: "内容", ChannelType: model.ChannelTypeEmail})

	o, err := svc.Overview()
	if err != nil {
		t.Fatalf("统计失败: %v", err)
	}
	if o.Channels != 1 || o.Templates != 1 || o.Topics != 1 || o.Recipients != 1 || o.Messages != 1 {
		t.Fatalf("统计计数不符: %+v", o)
	}
}

func TestRetryPolicyValidation(t *testing.T) {
	svc := newTestService()
	if _, err := svc.CreateRetryPolicy(model.RetryPolicy{Name: "策略", ChannelType: model.ChannelTypeEmail, MaxAttempts: 0}); err == nil {
		t.Fatalf("期望 MaxAttempts 必须大于 0 时报错")
	}
	if _, err := svc.CreateRetryPolicy(model.RetryPolicy{Name: "策略", ChannelType: model.ChannelTypeEmail, MaxAttempts: 3}); err != nil {
		t.Fatalf("创建策略失败: %v", err)
	}
}

func TestScheduleStateMachine(t *testing.T) {
	svc := newTestService()
	m, _ := svc.CreateMessage(model.Message{Title: "标题", Content: "内容", ChannelType: model.ChannelTypeEmail})
	sc, err := svc.CreateSchedule(model.Schedule{MessageID: m.ID, CronExpr: "0 9 * * *"})
	if err != nil {
		t.Fatalf("创建定时任务失败: %v", err)
	}
	if _, err := svc.ExecuteSchedule(sc.ID); err != nil {
		t.Fatalf("执行定时任务失败: %v", err)
	}
	if _, err := svc.CancelSchedule(sc.ID); err == nil {
		t.Fatalf("期望 executed 状态不可取消")
	}
}

func TestSendRecordStateMachine(t *testing.T) {
	svc := newTestService()
	m, _ := svc.CreateMessage(model.Message{Title: "标题", Content: "内容", ChannelType: model.ChannelTypeEmail})
	rp, _ := svc.CreateRecipient(model.Recipient{Name: "张三", ChannelType: model.ChannelTypeEmail, Address: "a@b.com"})
	rec, err := svc.CreateSendRecord(model.SendRecord{MessageID: m.ID, RecipientID: rp.ID, ChannelType: model.ChannelTypeEmail})
	if err != nil {
		t.Fatalf("创建发送记录失败: %v", err)
	}
	if _, err := svc.MarkRecordFailed(rec.ID, "网络错误"); err != nil {
		t.Fatalf("标记失败: %v", err)
	}
	if _, err := svc.RetryRecord(rec.ID); err != nil {
		t.Fatalf("重试失败: %v", err)
	}
	if _, err := svc.MarkRecordSuccess(rec.ID, 100); err != nil {
		t.Fatalf("标记成功失败: %v", err)
	}
}

func TestBatchDeleteTopics(t *testing.T) {
	svc := newTestService()
	tp1, _ := svc.CreateTopic(model.Topic{Name: "主题1"})
	tp2, _ := svc.CreateTopic(model.Topic{Name: "主题2"})
	n, err := svc.BatchDeleteTopics([]string{tp1.ID, tp2.ID})
	if err != nil {
		t.Fatalf("批量删除失败: %v", err)
	}
	if n != 2 {
		t.Fatalf("期望删除 2 个，得到 %d", n)
	}
}
