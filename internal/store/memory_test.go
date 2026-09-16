package store

import (
	"errors"
	"testing"
	"time"

	"notification/internal/model"
)

func newTestStore() *MemoryStore {
	return NewMemoryStore()
}

func mustChannel(t *testing.T, id, name, ctype string) *model.Channel {
	t.Helper()
	now := time.Now()
	return &model.Channel{ID: id, Name: name, Type: ctype, Status: model.ChannelEnabled, Config: "{}", Priority: 1, CreatedAt: now, UpdatedAt: now}
}

func TestChannelCRUD(t *testing.T) {
	s := newTestStore()
	c := mustChannel(t, "c1", "邮件渠道", model.ChannelTypeEmail)
	if err := s.CreateChannel(c); err != nil {
		t.Fatalf("创建渠道失败: %v", err)
	}
	if err := s.CreateChannel(mustChannel(t, "c2", "邮件渠道", model.ChannelTypeEmail)); !errors.Is(err, ErrConflict) {
		t.Fatalf("期望唯一性冲突，得到 %v", err)
	}
	got, err := s.GetChannel("c1")
	if err != nil || got.Name != "邮件渠道" {
		t.Fatalf("获取渠道失败: %v", err)
	}
	if len(s.ListChannels()) != 1 {
		t.Fatalf("期望 1 个渠道，得到 %d", len(s.ListChannels()))
	}
	if _, err := s.GetChannel("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
	got.Name = "邮件渠道2"
	if err := s.UpdateChannel(got); err != nil {
		t.Fatalf("更新渠道失败: %v", err)
	}
	if err := s.DeleteChannel("c1"); err != nil {
		t.Fatalf("删除渠道失败: %v", err)
	}
	if err := s.DeleteChannel("c1"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
}

func TestTemplateCRUD(t *testing.T) {
	s := newTestStore()
	tm := &model.Template{ID: "t1", Name: "欢迎模板", Type: model.ChannelTypeEmail, Subject: "欢迎", Content: "你好", Version: 1, Status: model.TemplateDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateTemplate(tm); err != nil {
		t.Fatalf("创建模板失败: %v", err)
	}
	if err := s.CreateTemplate(&model.Template{ID: "t2", Name: "欢迎模板", Type: model.ChannelTypeEmail, Version: 1}); !errors.Is(err, ErrConflict) {
		t.Fatalf("期望唯一性冲突，得到 %v", err)
	}
	if _, err := s.GetTemplateByName("欢迎模板"); err != nil {
		t.Fatalf("按名称获取模板失败: %v", err)
	}
	if err := s.DeleteTemplate("t1"); err != nil {
		t.Fatalf("删除模板失败: %v", err)
	}
}

func TestTopicCRUD(t *testing.T) {
	s := newTestStore()
	tp := &model.Topic{ID: "tp1", Name: "系统通知", Description: "系统通知主题", Status: model.TopicActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateTopic(tp); err != nil {
		t.Fatalf("创建主题失败: %v", err)
	}
	if err := s.CreateTopic(&model.Topic{ID: "tp2", Name: "系统通知"}); !errors.Is(err, ErrConflict) {
		t.Fatalf("期望唯一性冲突，得到 %v", err)
	}
	if _, err := s.GetTopic("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
	tp.Name = "系统通知2"
	if err := s.UpdateTopic(tp); err != nil {
		t.Fatalf("更新主题失败: %v", err)
	}
	if err := s.DeleteTopic("tp1"); err != nil {
		t.Fatalf("删除主题失败: %v", err)
	}
}

func TestRecipientCRUD(t *testing.T) {
	s := newTestStore()
	r := &model.Recipient{ID: "r1", Name: "张三", ChannelType: model.ChannelTypeEmail, Address: "a@b.com", Status: model.RecipientActive, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateRecipient(r); err != nil {
		t.Fatalf("创建接收人失败: %v", err)
	}
	if err := s.CreateRecipient(&model.Recipient{ID: "r2", Name: "李四", ChannelType: model.ChannelTypeEmail, Address: "a@b.com"}); !errors.Is(err, ErrConflict) {
		t.Fatalf("期望地址冲突，得到 %v", err)
	}
	if _, err := s.GetRecipientByAddress(model.ChannelTypeEmail, "a@b.com"); err != nil {
		t.Fatalf("按地址获取接收人失败: %v", err)
	}
	if err := s.DeleteRecipient("r1"); err != nil {
		t.Fatalf("删除接收人失败: %v", err)
	}
}

func TestSubscriptionCRUD(t *testing.T) {
	s := newTestStore()
	sub := &model.Subscription{ID: "s1", TopicID: "tp1", RecipientID: "r1", ChannelType: model.ChannelTypeEmail, Status: model.SubscriptionSubscribed, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateSubscription(sub); err != nil {
		t.Fatalf("创建订阅失败: %v", err)
	}
	if err := s.CreateSubscription(&model.Subscription{ID: "s2", TopicID: "tp1", RecipientID: "r1", ChannelType: model.ChannelTypeEmail}); !errors.Is(err, ErrConflict) {
		t.Fatalf("期望唯一订阅冲突，得到 %v", err)
	}
	if _, err := s.GetSubscriptionByTopicRecipient("tp1", "r1"); err != nil {
		t.Fatalf("按主题接收人获取订阅失败: %v", err)
	}
	if err := s.DeleteSubscription("s1"); err != nil {
		t.Fatalf("删除订阅失败: %v", err)
	}
}

func TestMessageCRUD(t *testing.T) {
	s := newTestStore()
	m := &model.Message{ID: "m1", Title: "标题", Content: "内容", ChannelType: model.ChannelTypeEmail, Priority: model.PriorityNormal, Status: model.MessageDraft, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateMessage(m); err != nil {
		t.Fatalf("创建消息失败: %v", err)
	}
	if _, err := s.GetMessage("m1"); err != nil {
		t.Fatalf("获取消息失败: %v", err)
	}
	if _, err := s.GetMessage("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
	m.Status = model.MessagePending
	if err := s.UpdateMessage(m); err != nil {
		t.Fatalf("更新消息失败: %v", err)
	}
	if err := s.DeleteMessage("m1"); err != nil {
		t.Fatalf("删除消息失败: %v", err)
	}
}

func TestSendRecordCRUD(t *testing.T) {
	s := newTestStore()
	rec := &model.SendRecord{ID: "sr1", MessageID: "m1", RecipientID: "r1", ChannelType: model.ChannelTypeEmail, Status: model.SendRecordPending, CreatedAt: time.Now()}
	if err := s.CreateSendRecord(rec); err != nil {
		t.Fatalf("创建发送记录失败: %v", err)
	}
	if _, err := s.GetSendRecord("sr1"); err != nil {
		t.Fatalf("获取发送记录失败: %v", err)
	}
	if _, err := s.GetSendRecord("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
	rec.Status = model.SendRecordSuccess
	if err := s.UpdateSendRecord(rec); err != nil {
		t.Fatalf("更新发送记录失败: %v", err)
	}
	if err := s.DeleteSendRecord("sr1"); err != nil {
		t.Fatalf("删除发送记录失败: %v", err)
	}
}

func TestRetryPolicyCRUD(t *testing.T) {
	s := newTestStore()
	p := &model.RetryPolicy{ID: "p1", Name: "邮件重试", ChannelType: model.ChannelTypeEmail, MaxAttempts: 3, BackoffMs: 1000, Status: model.RetryPolicyEnabled, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateRetryPolicy(p); err != nil {
		t.Fatalf("创建重试策略失败: %v", err)
	}
	if err := s.CreateRetryPolicy(&model.RetryPolicy{ID: "p2", Name: "邮件重试", ChannelType: model.ChannelTypeEmail}); !errors.Is(err, ErrConflict) {
		t.Fatalf("期望唯一性冲突，得到 %v", err)
	}
	if _, err := s.GetRetryPolicyByName("邮件重试"); err != nil {
		t.Fatalf("按名称获取策略失败: %v", err)
	}
	if err := s.DeleteRetryPolicy("p1"); err != nil {
		t.Fatalf("删除策略失败: %v", err)
	}
}

func TestScheduleCRUD(t *testing.T) {
	s := newTestStore()
	sc := &model.Schedule{ID: "sc1", MessageID: "m1", CronExpr: "0 9 * * *", Status: model.SchedulePending, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := s.CreateSchedule(sc); err != nil {
		t.Fatalf("创建定时任务失败: %v", err)
	}
	if _, err := s.GetSchedule("sc1"); err != nil {
		t.Fatalf("获取定时任务失败: %v", err)
	}
	if _, err := s.GetSchedule("missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("期望 ErrNotFound，得到 %v", err)
	}
	if err := s.DeleteSchedule("sc1"); err != nil {
		t.Fatalf("删除定时任务失败: %v", err)
	}
}
