# at-tray — Windows 系统托盘定时任务管理器

## 项目目标
纯 Go 写的 Windows 系统托盘工具，自带调度引擎（不用 at/schtasks），带 Web GUI 管理页面。

## 技术栈
- Go (pure Go, zero CGo)
- `github.com/gogpu/systray` — 系统托盘
- `net/http` — 本地 Web 服务器
- 嵌入式 HTML/CSS/JS 前端
- JSON 文件持久化 `%APPDATA%/at-tray/tasks.json`

## 已有文件
- `task.go` — Task 类型定义（已写好，别删字段）
- `storage.go` — JSON 持久化（已写好，别动 API）
- `main.go` — 入口（需要重写，当前是半成品）
- `PRD.md` — 需求文档
- `go.mod` — Go 模块
- `go.sum` — 依赖锁

## 核心需求

### 1. 系统托盘
托盘右键菜单：
- "打开管理页面" → 打开浏览器到本地 Web 界面
- "退出" → 关闭所有服务

### 2. Web GUI (硬编码在 Go 二进制里)
管理页面包含：
- **任务列表**：表格展示所有任务（时间/动作/重复/状态/操作按钮）
- **新建任务表单**：
  - 动作选择：关机 / 重启 / 锁定 / 命令
  - 如果选"命令"，显示文本输入框
  - 时间选择：快速预设按钮 [5分钟后] [15分钟] [30分钟] [1小时] [2小时] [4小时]
  - 自定义时间：HH:MM 输入 + 日期选择（今天/明天/日期输入）
  - 重复：单次（默认）/ 每天 / 每周
  - 最大次数：输入框（如果 MaxCount=1 则隐藏重复选项）
  - 提前通知：分钟数（0=不通知）
  - 重要通知：复选框（专注模式也显示）
  - 重启销毁：复选框（程序重启后不再保留）
  - 错过策略：跳过 / 立即执行
- **编辑/删除任务**

### 3. Web API
- `GET /api/tasks` — 列出所有任务
- `POST /api/tasks` — 创建任务
- `DELETE /api/tasks/{id}` — 删除任务
- `PATCH /api/tasks/{id}` — 修改任务（启用/禁用等）

### 4. 调度器
- 每 5 秒检查一次
- 到期执行：关机(shutdown /s /t 5)、重启(shutdown /r /t 5)、锁定(rundll32.exe user32.dll,LockWorkStation)、命令(exec.Command)
- 执行后：Executed++，检查是否达到 MaxCount → 自动禁用
- 提前通知：NotifMin 分钟前发气泡通知
- 重要通知：设置 Important 标记
- 错过策略：MissedSkip=可跳过，MissedExecute=立即执行
- 重启销毁：Persistent=false 的任务程序退出后不保存

### 5. 持久化
- 保存到 `%APPDATA%/at-tray/tasks.json`
- 只保存 Persistent=true 的任务
- 程序启动时加载

### 6. 通知
- 使用 tray.ShowNotification()
- Important=true 时用特别的 Windows 通知（高优先级）

## 构建
```bash
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o at-tray.exe
```

## 前端设计原则
- 简洁、深色主题
- 所有代码在单个 HTML 文件里（CSS + JS 内联）
- 无外部 CDN 依赖
- 时间选择器要有快速预设按钮
- 响应式，窗口大小自适应
