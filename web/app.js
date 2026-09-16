// 通知推送中心前端逻辑
const API_KEY = 'dev-notify-key';

async function api(path) {
  const res = await fetch(path, { headers: { 'X-API-Key': API_KEY } });
  if (!res.ok) {
    throw new Error('请求失败: ' + res.status);
  }
  return res.json();
}

async function loadOverview() {
  const body = await api('/api/stats/overview');
  const data = body.data || {};
  const cards = [
    { label: '消息总数', value: data.messages },
    { label: '已发送', value: data.messages_sent },
    { label: '待发送', value: data.messages_pending },
    { label: '发送失败', value: data.messages_failed },
    { label: '渠道数', value: data.channels },
    { label: '模板数', value: data.templates },
    { label: '主题数', value: data.topics },
    { label: '接收人', value: data.recipients },
    { label: '订阅数', value: data.total_subscribers },
    { label: '发送记录', value: data.send_records },
    { label: '成功记录', value: data.records_succeeded },
    { label: '成功率', value: (data.success_rate || 0).toFixed(1) + '%' },
  ];
  const container = document.getElementById('stats');
  container.innerHTML = '';
  cards.forEach(function (c) {
    const div = document.createElement('div');
    div.className = 'stat-card';
    div.innerHTML = '<div class="label">' + c.label + '</div><div class="value">' + c.value + '</div>';
    container.appendChild(div);
  });
}

function badgeClass(status) {
  const map = {
    sent: 'sent', pending: 'pending', draft: 'draft',
    sending: 'sending', failed: 'failed', cancelled: 'cancelled',
  };
  return map[status] || '';
}

async function loadMessages() {
  const body = await api('/api/messages?size=50');
  const items = (body.data && body.data.items) || [];
  const tbody = document.querySelector('#message-table tbody');
  tbody.innerHTML = '';
  items.forEach(function (m) {
    const tr = document.createElement('tr');
    const created = new Date(m.created_at).toLocaleString();
    tr.innerHTML =
      '<td>' + m.id + '</td>' +
      '<td>' + m.title + '</td>' +
      '<td>' + m.channel_type + '</td>' +
      '<td>' + m.priority + '</td>' +
      '<td><span class="badge ' + badgeClass(m.status) + '">' + m.status + '</span></td>' +
      '<td>' + created + '</td>';
    tbody.appendChild(tr);
  });
}

async function loadChannels() {
  const body = await api('/api/channels?size=50');
  const items = (body.data && body.data.items) || [];
  const tbody = document.querySelector('#channel-table tbody');
  tbody.innerHTML = '';
  items.forEach(function (c) {
    const tr = document.createElement('tr');
    tr.innerHTML =
      '<td>' + c.name + '</td>' +
      '<td>' + c.type + '</td>' +
      '<td>' + c.status + '</td>' +
      '<td>' + c.priority + '</td>';
    tbody.appendChild(tr);
  });
}

async function load() {
  await loadOverview();
  await loadMessages();
  await loadChannels();
}

document.getElementById('refresh-btn').addEventListener('click', load);
load();
