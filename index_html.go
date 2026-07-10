package main

const indexHTML = `<!DOCTYPE html>
<html lang="zh-CN">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width,initial-scale=1">
<title>at-tray 定时任务管理器</title>
<style>
:root {
  --bg: #1a1a2e;
  --surface: #16213e;
  --surface2: #1e2d50;
  --border: #2a3a6a;
  --text: #e0e0e0;
  --text-dim: #8899bb;
  --primary: #4a8eff;
  --primary-hover: #3a7aee;
  --danger: #e74c5c;
  --danger-hover: #c0392b;
  --success: #2ecc71;
  --warning: #f39c12;
  --radius: 8px;
  --shadow: 0 4px 24px rgba(0,0,0,0.4);
}
* { margin: 0; padding: 0; box-sizing: border-box; }
body {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", "Microsoft YaHei", sans-serif;
  background: var(--bg);
  color: var(--text);
  line-height: 1.6;
  min-height: 100vh;
}
header {
  background: var(--surface);
  padding: 16px 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border);
  flex-wrap: wrap;
  gap: 12px;
}
header h1 { font-size: 20px; font-weight: 600; }
.header-actions { display: flex; gap: 10px; align-items: center; }
.btn {
  padding: 8px 18px;
  border: none;
  border-radius: var(--radius);
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
  transition: background .15s, transform .1s;
  display: inline-flex;
  align-items: center;
  gap: 6px;
}
.btn:active { transform: scale(.97); }
.btn-primary { background: var(--primary); color: #fff; }
.btn-primary:hover { background: var(--primary-hover); }
.btn-danger { background: var(--danger); color: #fff; }
.btn-danger:hover { background: var(--danger-hover); }
.btn-ghost {
  background: transparent;
  color: var(--text-dim);
  border: 1px solid var(--border);
}
.btn-ghost:hover { background: var(--surface2); color: var(--text); }
.btn-sm { padding: 4px 12px; font-size: 12px; }

/* Main content */
main {
  padding: 24px;
  max-width: 1100px;
  margin: 0 auto;
}

/* Task table */
.table-wrap {
  background: var(--surface);
  border-radius: var(--radius);
  border: 1px solid var(--border);
  overflow: hidden;
}
table {
  width: 100%;
  border-collapse: collapse;
}
th {
  text-align: left;
  padding: 12px 16px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-dim);
  text-transform: uppercase;
  letter-spacing: .5px;
  border-bottom: 1px solid var(--border);
  background: var(--surface2);
}
td {
  padding: 12px 16px;
  font-size: 14px;
  border-bottom: 1px solid var(--border);
  vertical-align: middle;
}
tr:last-child td { border-bottom: none; }
tr:hover td { background: rgba(74,142,255,.05); }
td.actions { white-space: nowrap; }

/* Badge / tag */
.badge {
  display: inline-block;
  padding: 2px 10px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
}
.badge-active { background: rgba(46,204,113,.15); color: var(--success); }
.badge-disabled { background: rgba(153,153,153,.15); color: #999; }
.badge-done { background: rgba(243,156,18,.15); color: var(--warning); }
.badge-shutdown { background: rgba(231,76,92,.15); color: var(--danger); }
.badge-restart { background: rgba(243,156,18,.15); color: var(--warning); }
.badge-lock { background: rgba(52,152,219,.15); color: #3498db; }
.badge-command { background: rgba(155,89,182,.15); color: #9b59b6; }
.badge-once { background: rgba(74,142,255,.15); color: var(--primary); }
.badge-daily { background: rgba(46,204,113,.15); color: var(--success); }
.badge-weekly { background: rgba(243,156,18,.15); color: var(--warning); }
.badge-hourly { background: rgba(52,152,219,.15); color: #3498db; }

/* Empty state */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  color: var(--text-dim);
}
.empty-state p { font-size: 16px; margin-bottom: 16px; }

/* Modal */
.modal-overlay {
  position: fixed;
  top: 0; left: 0; right: 0; bottom: 0;
  background: rgba(0,0,0,.6);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}
.modal-overlay.hidden { display: none; }
.modal-content {
  background: var(--surface);
  border-radius: var(--radius);
  border: 1px solid var(--border);
  box-shadow: var(--shadow);
  width: 100%;
  max-width: 580px;
  max-height: 90vh;
  overflow-y: auto;
  padding: 24px;
}
.modal-content h2 {
  font-size: 18px;
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid var(--border);
}

/* Form */
.form-group {
  margin-bottom: 16px;
}
.form-group label {
  display: block;
  font-size: 13px;
  font-weight: 500;
  color: var(--text-dim);
  margin-bottom: 6px;
}
.form-group select,
.form-group input[type="text"],
.form-group input[type="number"],
.form-group input[type="date"],
.form-group input[type="time"] {
  width: 100%;
  padding: 8px 12px;
  background: var(--bg);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  color: var(--text);
  font-size: 14px;
  outline: none;
  transition: border-color .15s;
}
.form-group select:focus,
.form-group input:focus { border-color: var(--primary); }
.form-group select { cursor: pointer; }
.form-row {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}
.form-row .form-group { flex: 1; min-width: 120px; }
.form-inline {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}
.form-inline label { margin-bottom: 0; font-size: 14px; }
.form-inline input[type="checkbox"] {
  width: 16px;
  height: 16px;
  accent-color: var(--primary);
}
.checkbox-group {
  display: flex;
  align-items: center;
  gap: 8px;
}
.checkbox-group input[type="checkbox"] { width: 16px; height: 16px; accent-color: var(--primary); }
.checkbox-group label { margin-bottom: 0; cursor: pointer; }

/* Radio group */
.radio-group {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;
}
.radio-group label {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 14px;
  font-weight: 400;
  color: var(--text);
  cursor: pointer;
  margin-bottom: 0;
}
.radio-group input[type="radio"] { accent-color: var(--primary); width: 16px; height: 16px; }

/* Time presets */
.time-presets {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}
.time-presets button {
  padding: 6px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius);
  background: var(--bg);
  color: var(--text);
  cursor: pointer;
  font-size: 13px;
  transition: background .15s, border-color .15s;
}
.time-presets button:hover {
  background: var(--surface2);
  border-color: var(--primary);
}
.time-presets button.active {
  background: var(--primary);
  border-color: var(--primary);
  color: #fff;
}

/* Modal footer */
.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 20px;
  padding-top: 16px;
  border-top: 1px solid var(--border);
}

/* Hidden */
.hidden { display: none !important; }

/* Responsive */
@media(max-width:600px) {
  header { padding: 12px 16px; }
  header h1 { font-size: 16px; }
  main { padding: 12px; }
  th, td { padding: 8px 10px; font-size: 13px; }
  .modal-content { padding: 16px; }
  .form-row .form-group { min-width: 100%; }
  .radio-group { gap: 10px; }
}

/* Scrollbar */
::-webkit-scrollbar { width: 6px; }
::-webkit-scrollbar-track { background: var(--bg); }
::-webkit-scrollbar-thumb { background: var(--border); border-radius: 3px; }
</style>
</head>
<body>

<header>
  <h1>⏱ at-tray</h1>
  <div class="header-actions">
    <button class="btn btn-primary" id="newTaskBtn">+ 新建任务</button>
  </div>
</header>

<main>
  <div class="table-wrap">
    <table>
      <thead>
        <tr>
          <th>时间</th>
          <th>动作</th>
          <th>重复</th>
          <th>执行</th>
          <th>状态</th>
          <th>操作</th>
        </tr>
      </thead>
      <tbody id="taskList"></tbody>
    </table>
    <div class="empty-state" id="emptyState">
      <p>暂无定时任务</p>
      <button class="btn btn-primary" onclick="openNewTask()">+ 创建第一个任务</button>
    </div>
  </div>
</main>

<!-- Modal -->
<div class="modal-overlay hidden" id="modalOverlay">
  <div class="modal-content">
    <h2 id="modalTitle">新建任务</h2>
    <form id="taskForm" onsubmit="return false">
      <!-- 动作 -->
      <div class="form-group">
        <label for="formAction">动作</label>
        <select id="formAction" onchange="onActionChange()">
          <option value="0">🛑 关机</option>
          <option value="1">🔄 重启</option>
          <option value="2">🔒 锁定</option>
          <option value="3">💻 命令</option>
        </select>
      </div>

      <!-- 命令输入 -->
      <div class="form-group hidden" id="commandGroup">
        <label for="formCommand">命令内容</label>
        <input type="text" id="formCommand" placeholder="例如: notepad.exe">
      </div>

      <!-- 时间模式 -->
      <div class="form-group">
        <label>时间模式</label>
        <div class="radio-group">
          <label><input type="radio" name="timeMode" value="absolute" checked> 绝对时间</label>
          <label><input type="radio" name="timeMode" value="relative"> 相对时间</label>
        </div>
      </div>

      <!-- 相对时间 -->
      <div class="form-group hidden" id="relativeTimeGroup">
        <label>相对时间</label>
        <div class="form-inline">
          <input type="number" id="formRelativeValue" min="1" value="5" style="width:80px">
          <select id="formRelativeUnit">
            <option value="minutes">分钟</option>
            <option value="hours">小时</option>
          </select>
          <span>后执行</span>
        </div>
      </div>

      <!-- 时间 -->
      <div class="form-group" id="absoluteTimeGroup">
        <label>执行时间</label>
        <div class="time-presets" id="timePresets">
          <button data-min="5">5 分钟后</button>
          <button data-min="15">15 分钟后</button>
          <button data-min="30">30 分钟后</button>
          <button data-min="60">1 小时后</button>
          <button data-min="120">2 小时后</button>
          <button data-min="240">4 小时后</button>
        </div>
        <div class="form-row">
          <div style="flex:1;min-width:100px">
            <input type="time" id="formTime" value="09:00">
          </div>
          <div style="flex:1;min-width:120px">
            <select id="formDateMode" onchange="onDateModeChange()">
              <option value="today">今天</option>
              <option value="tomorrow">明天</option>
              <option value="custom">指定日期</option>
            </select>
          </div>
          <div style="flex:1;min-width:130px" id="customDateGroup" class="hidden">
            <input type="date" id="formDate">
          </div>
        </div>
      </div>

      <!-- 重复 -->
      <div class="form-group hidden" id="repeatGroup">
        <label>重复方式</label>
        <div class="radio-group">
          <label><input type="radio" name="repeat" value="1" checked> 每天</label>
          <label><input type="radio" name="repeat" value="2"> 每周</label>
          <label><input type="radio" name="repeat" value="3"> 每小时</label>
        </div>
      </div>

      <!-- 最大次数 -->
      <div class="form-group">
        <label for="formMaxCount">最大执行次数</label>
        <input type="number" id="formMaxCount" min="0" value="1" style="width:120px">
        <div style="font-size:12px;color:var(--text-dim);margin-top:4px">设置为 0 表示无限制</div>
      </div>

      <!-- 提前通知 -->
      <div class="form-group">
        <label for="formNotifyMin">提前通知（分钟，0=不通知）</label>
        <input type="number" id="formNotifyMin" min="0" value="0" style="width:120px">
      </div>

      <!-- 选项 -->
      <div class="form-row">
        <div class="form-group">
          <div class="checkbox-group">
            <input type="checkbox" id="formImportant">
            <label for="formImportant">重要通知（专注模式也显示）</label>
          </div>
        </div>
        <div class="form-group">
          <div class="checkbox-group">
            <input type="checkbox" id="formPersistent" checked>
            <label for="formPersistent">持久化保留（重启后不销毁）</label>
          </div>
        </div>
      </div>

      <!-- 错过策略 -->
      <div class="form-group">
        <label>错过策略</label>
        <div class="radio-group">
          <label><input type="radio" name="missed" value="0" checked> 跳过</label>
          <label><input type="radio" name="missed" value="1"> 立即执行</label>
        </div>
      </div>
      <!-- 高级选项 -->
      <div class="form-group">
        <details style="cursor:pointer">
          <summary style="color:var(--text-dim);font-size:13px">⚙ 高级选项</summary>
          <div style="margin-top:10px">
            <label for="formCron" style="font-size:13px;color:var(--text-dim);display:block;margin-bottom:4px">Cron 表达式（预留）</label>
            <input type="text" id="formCron" placeholder="*/5 * * * *" style="width:100%;padding:8px 12px;background:var(--bg);border:1px solid var(--border);border-radius:var(--radius);color:var(--text);font-size:14px;outline:none">
            <div style="font-size:12px;color:var(--text-dim);margin-top:4px">例如: */5 * * * *（每5分钟执行一次）</div>
          </div>
        </details>
      </div>
    </form>
    <div class="modal-footer">
      <button class="btn btn-ghost" id="cancelBtn">取消</button>
      <button class="btn btn-primary" id="saveBtn">保存</button>
    </div>
  </div>
</div>

<script>
// ── State ──
let editingId = null;

// ── API ──
async function api(method, path, body) {
  const opts = { method, headers: { 'Content-Type': 'application/json' } };
  if (body !== undefined) opts.body = JSON.stringify(body);
  const res = await fetch(path, opts);
  if (!res.ok) throw new Error(await res.text());
  return res.json();
}

// ── Format helpers ──
const ACTIONS = ['关 机', '重 启', '锁 定', '命 令'];
function actionBadge(a) {
  const cls = ['badge-shutdown','badge-restart','badge-lock','badge-command'][a]||'';
  return '<span class="badge '+cls+'">'+ACTIONS[a]+'</span>';
}
const REPEATS = ['每天','每周','每小时'];
function repeatBadge(r) {
  const cls = ['badge-once','badge-daily','badge-weekly','badge-hourly'][r]||'';
  return '<span class="badge '+cls+'">'+REPEATS[r]+'</span>';
}
function statusBadge(t) {
  if (!t.enabled) return '<span class="badge badge-disabled">已禁用</span>';
  if (t.max_count>0 && t.executed>=t.max_count) return '<span class="badge badge-done">已完成</span>';
  return '<span class="badge badge-active">运行中</span>';
}
function fmtTime(iso) {
  const d = new Date(iso);
  const pad = n => String(n).padStart(2,'0');
  return d.getFullYear()+'-'+pad(d.getMonth()+1)+'-'+pad(d.getDate())+' '+pad(d.getHours())+':'+pad(d.getMinutes());
}

// ── Render task list ──
async function loadTasks() {
  const tasks = await api('GET','/api/tasks');
  const tbody = document.getElementById('taskList');
  const empty = document.getElementById('emptyState');
  if (!tasks || tasks.length===0) {
    tbody.innerHTML = '';
    empty.style.display = 'block';
    return;
  }
  empty.style.display = 'none';
  tbody.innerHTML = tasks.map(t => {
    const execInfo = t.max_count>0 ? t.executed+'/'+t.max_count : t.executed+'/无限';
    return '<tr>' +
      '<td>'+fmtTime(t.target_time)+'</td>' +
      '<td>'+actionBadge(t.action)+(t.action===3&&t.command?'<br><small style="color:var(--text-dim)">'+escHtml(t.command)+'</small>':'')+'</td>' +
      '<td>'+repeatBadge(t.repeat)+'</td>' +
      '<td>'+execInfo+'</td>' +
      '<td>'+statusBadge(t)+'</td>' +
      '<td class="actions">' +
        '<button class="btn btn-sm btn-ghost" onclick="editTask(\''+t.id+'\')">✏️</button> ' +
        '<button class="btn btn-sm btn-ghost" onclick="toggleTask(\''+t.id+'\','+t.enabled+')">'+(t.enabled?'⏸':'▶️')+'</button> ' +
        '<button class="btn btn-sm btn-danger" onclick="deleteTask(\''+t.id+'\')">🗑</button>' +
      '</td></tr>';
  }).join('');
}

function escHtml(s) {
  const d = document.createElement('div');
  d.textContent = s;
  return d.innerHTML;
}

// ── CRUD ──
async function deleteTask(id) {
  if (!confirm('确定删除此任务？')) return;
  await api('DELETE','/api/tasks/'+id);
  loadTasks();
}

async function toggleTask(id, current) {
  // 重新启用时重置执行次数
  await api('PATCH','/api/tasks/'+id, { enabled: !current, executed: current ? undefined : 0 });
  loadTasks();
}

// ── Modal ──
function openNewTask() {
  editingId = null;
  document.getElementById('modalTitle').textContent = '新建任务';
  document.getElementById('taskForm').reset();
  document.getElementById('formMaxCount').value = 1;
  document.getElementById('formNotifyMin').value = 0;
  document.getElementById('formPersistent').checked = true;
  document.getElementById('formImportant').checked = false;
  document.getElementById('formDateMode').value = 'today';
  document.getElementById('formDate').value = todayStr();
  document.getElementById('customDateGroup').classList.add('hidden');
  document.querySelector('input[name="repeat"][value="1"]').checked = true;
  document.querySelector('input[name="missed"][value="0"]').checked = true;
  // Reset time mode to absolute
  document.querySelector('input[name="timeMode"][value="absolute"]').checked = true;
  document.getElementById('absoluteTimeGroup').classList.remove('hidden');
  document.getElementById('relativeTimeGroup').classList.add('hidden');
  document.getElementById('formAction').value = '0';
  onActionChange();
  if (parseInt(document.getElementById('formMaxCount').value) !== 1) {
    document.getElementById('repeatGroup').classList.remove('hidden');
  }
  document.getElementById('commandGroup').classList.add('hidden');
  clearPresetActive();
  // Set default time to next rounded 5 min
  const now = new Date();
  const min = Math.ceil(now.getMinutes()/5)*5;
  now.setMinutes(min, 0, 0);
  document.getElementById('formTime').value = pad2(now.getHours())+':'+pad2(now.getMinutes()%60);
  showModal();
}

async function editTask(id) {
  const tasks = await api('GET','/api/tasks');
  const t = tasks.find(x => x.id===id);
  if (!t) return;
  editingId = id;
  document.getElementById('modalTitle').textContent = '编辑任务';

  document.getElementById('formAction').value = String(t.action);
  onActionChange();
  if (t.action===3) document.getElementById('formCommand').value = t.command||'';

  const d = new Date(t.target_time);
  document.getElementById('formTime').value = pad2(d.getHours())+':'+pad2(d.getMinutes());
  const today = todayStr();
  const tomorrow = tomorrowStr();
  const dateStr = d.getFullYear()+'-'+pad2(d.getMonth()+1)+'-'+pad2(d.getDate());
  if (dateStr===today) {
    document.getElementById('formDateMode').value = 'today';
    document.getElementById('customDateGroup').classList.add('hidden');
  } else if (dateStr===tomorrow) {
    document.getElementById('formDateMode').value = 'tomorrow';
    document.getElementById('customDateGroup').classList.add('hidden');
  } else {
    document.getElementById('formDateMode').value = 'custom';
    document.getElementById('customDateGroup').classList.remove('hidden');
    document.getElementById('formDate').value = dateStr;
  }

  const repRadio = document.querySelector('input[name="repeat"][value="'+t.repeat+'"]');
  if (repRadio) repRadio.checked = true;
  if (t.repeat>0) document.getElementById('repeatGroup').classList.remove('hidden');
  document.getElementById('formMaxCount').value = t.max_count;
  if (t.max_count===1) document.getElementById('repeatGroup').classList.add('hidden');

  document.getElementById('formNotifyMin').value = t.notify_min;
  document.getElementById('formImportant').checked = t.important;
  document.getElementById('formPersistent').checked = t.persistent;
  document.querySelector('input[name="missed"][value="'+t.missed_policy+'"]').checked = true;

  // 编辑时强制使用绝对时间模式
  document.querySelector('input[name="timeMode"][value="absolute"]').checked = true;
  document.getElementById('absoluteTimeGroup').classList.remove('hidden');
  document.getElementById('relativeTimeGroup').classList.add('hidden');

  clearPresetActive();
  showModal();
}

function showModal() {
  document.getElementById('modalOverlay').classList.remove('hidden');
}

function hideModal() {
  document.getElementById('modalOverlay').classList.add('hidden');
}

// ── Form helpers ──
function pad2(n) { return String(n).padStart(2,'0'); }
function todayStr() {
  const d = new Date();
  return d.getFullYear()+'-'+pad2(d.getMonth()+1)+'-'+pad2(d.getDate());
}
function tomorrowStr() {
  const d = new Date();
  d.setDate(d.getDate()+1);
  return d.getFullYear()+'-'+pad2(d.getMonth()+1)+'-'+pad2(d.getDate());
}

function onActionChange() {
  const v = parseInt(document.getElementById('formAction').value);
  document.getElementById('commandGroup').classList.toggle('hidden', v!==3);
}

function onDateModeChange() {
  const v = document.getElementById('formDateMode').value;
  document.getElementById('customDateGroup').classList.toggle('hidden', v!=='custom');
}

// 时间模式切换
document.querySelectorAll('input[name="timeMode"]').forEach(el => {
  el.addEventListener('change', function() {
    const isAbsolute = this.value === 'absolute';
    document.getElementById('absoluteTimeGroup').classList.toggle('hidden', !isAbsolute);
    document.getElementById('relativeTimeGroup').classList.toggle('hidden', isAbsolute);
  });
});

function clearPresetActive() {
  document.querySelectorAll('#timePresets button').forEach(b => b.classList.remove('active'));
}

// Time presets
document.getElementById('timePresets').addEventListener('click', function(e) {
  if (e.target.tagName!=='BUTTON') return;
  clearPresetActive();
  e.target.classList.add('active');
  const min = parseInt(e.target.dataset.min);
  const now = new Date();
  now.setMinutes(now.getMinutes()+min);
  document.getElementById('formTime').value = pad2(now.getHours())+':'+pad2(now.getMinutes()%60);
  // Set date
  const today = todayStr();
  const tom = tomorrowStr();
  const dStr = now.getFullYear()+'-'+pad2(now.getMonth()+1)+'-'+pad2(now.getDate());
  if (dStr===today) {
    document.getElementById('formDateMode').value = 'today';
    document.getElementById('customDateGroup').classList.add('hidden');
  } else if (dStr===tom) {
    document.getElementById('formDateMode').value = 'tomorrow';
    document.getElementById('customDateGroup').classList.add('hidden');
  } else {
    document.getElementById('formDateMode').value = 'custom';
    document.getElementById('customDateGroup').classList.remove('hidden');
    document.getElementById('formDate').value = dStr;
  }
});

// MaxCount change → hide repeat when 1
document.getElementById('formMaxCount').addEventListener('input', function() {
  const v = parseInt(this.value)||1;
  document.getElementById('repeatGroup').classList.toggle('hidden', v===1);
});

// ── Save ──
document.getElementById('saveBtn').addEventListener('click', async function() {
  const action = parseInt(document.getElementById('formAction').value);
  const command = document.getElementById('formCommand').value;
  const timeMode = document.querySelector('input[name="timeMode"]:checked').value;
  let targetTime;
  if (timeMode === 'relative') {
    const relValue = parseInt(document.getElementById('formRelativeValue').value) || 0;
    const relUnit = document.getElementById('formRelativeUnit').value;
    if (relValue <= 0) { alert('请输入有效的时间'); return; }
    targetTime = new Date();
    if (relUnit === 'minutes') {
      targetTime.setMinutes(targetTime.getMinutes() + relValue);
    } else {
      targetTime.setHours(targetTime.getHours() + relValue);
    }
    targetTime.setSeconds(0, 0);
  } else {
    const timeStr = document.getElementById('formTime').value;
    if (!timeStr) { alert('请选择时间'); return; }
    const dateMode = document.getElementById('formDateMode').value;
    let dateStr;
    if (dateMode==='today') {
      const d = new Date(); dateStr = d.getFullYear()+'-'+pad2(d.getMonth()+1)+'-'+pad2(d.getDate());
    } else if (dateMode==='tomorrow') {
      const d = new Date(); d.setDate(d.getDate()+1); dateStr = d.getFullYear()+'-'+pad2(d.getMonth()+1)+'-'+pad2(d.getDate());
    } else {
      dateStr = document.getElementById('formDate').value;
      if (!dateStr) { alert('请选择日期'); return; }
    }
    targetTime = new Date(dateStr+'T'+timeStr+':00');
  }
  if (isNaN(targetTime.getTime())) { alert('时间格式无效'); return; }

  const repeat = parseInt(document.querySelector('input[name="repeat"]:checked').value);
  const maxCount = parseInt(document.getElementById('formMaxCount').value)||0;
  const notifyMin = parseInt(document.getElementById('formNotifyMin').value)||0;
  const important = document.getElementById('formImportant').checked;
  const persistent = document.getElementById('formPersistent').checked;
  const missedPolicy = parseInt(document.querySelector('input[name="missed"]:checked').value);

  if (action===3 && !command.trim()) { alert('请输入命令内容'); return; }

  const payload = {
    action,
    command: action===3 ? command : '',
    target_time: targetTime.toISOString(),
    repeat,
    max_count: maxCount===0 ? 0 : maxCount,
    notify_min: notifyMin,
    important,
    persistent,
    missed_policy: missedPolicy
  };

  // Validate max count: repeat tasks need a positive max count
  if (payload.max_count===0 && repeat>0) payload.max_count = 0;

  try {
    if (editingId) {
      // 编辑时重置执行次数
      payload.executed = 0;
      await api('PATCH','/api/tasks/'+editingId, payload);
    } else {
      await api('POST','/api/tasks', payload);
    }
    hideModal();
    loadTasks();
  } catch(e) {
    alert('保存失败: '+e.message);
  }
});

// ── Cancel ──
document.getElementById('cancelBtn').addEventListener('click', hideModal);
document.getElementById('modalOverlay').addEventListener('click', function(e) {
  if (e.target===this) hideModal();
});

// ── New task button ──
document.getElementById('newTaskBtn').addEventListener('click', openNewTask);

// Keyboard: Escape to close
document.addEventListener('keydown', function(e) {
  if (e.key==='Escape') hideModal();
});

// ── Init ──
loadTasks();
// Refresh list every 10s
setInterval(loadTasks, 10000);
</script>
</body>
</html>`
