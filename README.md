# 通知推送中心（notification）

纯 Go 标准库实现的通知推送中心后端服务，零第三方依赖（仅 `net/http` + 标准库）。覆盖消息推送、订阅管理、定时发送与发送记录追踪。

## 运行

```bash
# 使用本机 Go 1.22 工具链
/Users/fengyin/.local/go/bin/go run ./cmd/server

# 或构建后运行
/Users/fengyin/.local/go/bin/go build -o bin/server ./cmd/server
./bin/server
```

环境变量（均有默认值）：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| `ADDR` / `PORT` | `:8080` | 监听地址 |
| `API_KEY` | `dev-notify-key` | API 鉴权密钥（写请求头 `X-API-Key`） |
| `MAX_PAGE_SIZE` | `100` | 分页单页上限 |
| `RATE_LIMIT_RPS` | `100` | 令牌桶每秒补充令牌数 |
| `RATE_BURST` | `200` | 令牌桶容量 |
| `LOG_LEVEL` | `info` | 日志级别（debug/info/warn/error） |

服务启动后可通过 HTTP API 管理通知配置和发送流程。

所有 `/api/*` 请求需携带请求头 `X-API-Key: dev-notify-key`。

## API 一览

统一响应结构：`{"code":0,"message":"ok","data":...}`；错误码：400 校验失败、401 鉴权失败、404 不存在、409 冲突、429 限流、500 服务器错误。

### 渠道 Channel

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/channels` | 创建渠道 |
| GET | `/api/channels` | 分页查询（`type`/`status`/`keyword` 筛选） |
| GET | `/api/channels/{id}` | 获取渠道 |
| PUT | `/api/channels/{id}` | 更新渠道 |
| DELETE | `/api/channels/{id}` | 删除渠道 |
| POST | `/api/channels/batch-delete` | 批量删除（`{"ids":[...]}`） |

### 模板 Template

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/templates` | 创建模板 |
| GET | `/api/templates` | 分页查询（`type`/`status`/`keyword`） |
| GET | `/api/templates/{id}` | 获取模板 |
| PUT | `/api/templates/{id}` | 更新模板（版本号 +1） |
| DELETE | `/api/templates/{id}` | 删除模板 |
| POST | `/api/templates/{id}/transition` | 状态流转（draft→active→archived） |

### 主题 Topic

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/topics` | 创建主题 |
| GET | `/api/topics` | 分页查询（`status`/`keyword`） |
| GET | `/api/topics/{id}` | 获取主题 |
| PUT | `/api/topics/{id}` | 更新主题 |
| DELETE | `/api/topics/{id}` | 删除主题 |
| POST | `/api/topics/batch-delete` | 批量删除 |

### 接收人 Recipient

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/recipients` | 创建接收人 |
| GET | `/api/recipients` | 分页查询（`channel_type`/`group`/`status`/`keyword`） |
| GET | `/api/recipients/{id}` | 获取接收人 |
| PUT | `/api/recipients/{id}` | 更新接收人 |
| DELETE | `/api/recipients/{id}` | 删除接收人 |
| POST | `/api/recipients/{id}/status` | 更新状态（active/unsubscribed） |
| POST | `/api/recipients/batch-delete` | 批量删除 |

### 订阅 Subscription

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/subscriptions` | 创建订阅（校验主题与接收人存在） |
| GET | `/api/subscriptions` | 分页查询（`topic_id`/`recipient_id`/`channel_type`/`status`） |
| GET | `/api/subscriptions/{id}` | 获取订阅 |
| PUT | `/api/subscriptions/{id}` | 更新订阅 |
| DELETE | `/api/subscriptions/{id}` | 删除订阅 |
| POST | `/api/subscriptions/{id}/unsubscribe` | 退订 |
| POST | `/api/subscriptions/{id}/subscribe` | 重新订阅 |
| POST | `/api/subscriptions/batch-delete` | 批量删除 |

### 消息 Message

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/messages` | 创建消息（校验模板与主题存在） |
| GET | `/api/messages` | 分页查询（`template_id`/`topic_id`/`channel_type`/`priority`/`status`/`keyword`） |
| GET | `/api/messages/{id}` | 获取消息 |
| PUT | `/api/messages/{id}` | 更新消息（仅 draft/pending） |
| DELETE | `/api/messages/{id}` | 删除消息 |
| POST | `/api/messages/{id}/transition` | 状态流转 |
| POST | `/api/messages/{id}/send` | 发送消息（生成发送记录） |
| POST | `/api/messages/batch-send` | 批量发送 |
| POST | `/api/messages/batch-status` | 批量更新状态 |
| POST | `/api/messages/batch-delete` | 批量删除 |

消息状态机：`draft→pending→sending→sent/failed`；`pending→cancelled`。

### 发送记录 SendRecord

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/send-records` | 创建发送记录 |
| GET | `/api/send-records` | 分页查询（`message_id`/`recipient_id`/`channel_type`/`status`） |
| GET | `/api/send-records/{id}` | 获取发送记录 |
| PUT | `/api/send-records/{id}` | 更新发送记录 |
| DELETE | `/api/send-records/{id}` | 删除发送记录 |
| POST | `/api/send-records/{id}/success` | 标记成功 |
| POST | `/api/send-records/{id}/failed` | 标记失败 |
| POST | `/api/send-records/{id}/retry` | 重试 |
| POST | `/api/send-records/batch-status` | 批量更新状态 |

发送记录状态机：`pending→success/failed`；`failed→retrying→success/failed`。

### 重试策略 RetryPolicy

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/retry-policies` | 创建重试策略 |
| GET | `/api/retry-policies` | 分页查询（`channel_type`/`status`/`keyword`） |
| GET | `/api/retry-policies/{id}` | 获取策略 |
| PUT | `/api/retry-policies/{id}` | 更新策略 |
| DELETE | `/api/retry-policies/{id}` | 删除策略 |
| POST | `/api/retry-policies/batch-delete` | 批量删除 |

### 定时任务 Schedule

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/schedules` | 创建定时任务（校验消息存在） |
| GET | `/api/schedules` | 分页查询（`message_id`/`status`） |
| GET | `/api/schedules/{id}` | 获取任务 |
| PUT | `/api/schedules/{id}` | 更新任务 |
| DELETE | `/api/schedules/{id}` | 删除任务 |
| POST | `/api/schedules/{id}/execute` | 执行（pending→executed） |
| POST | `/api/schedules/{id}/cancel` | 取消（pending→cancelled） |
| POST | `/api/schedules/batch-delete` | 批量删除 |

### 数据导出

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/export` | 导出全部数据汇总快照 JSON |

## 测试

```bash
/Users/fengyin/.local/go/bin/go test ./...
```
