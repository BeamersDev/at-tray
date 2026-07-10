package main

import "time"

type ActionType int

const (
	ActionShutdown ActionType = iota
	ActionRestart
	ActionLock
	ActionCommand
)

func (a ActionType) String() string {
	switch a {
	case ActionShutdown:
		return "shutdown"
	case ActionRestart:
		return "restart"
	case ActionLock:
		return "lock"
	case ActionCommand:
		return "command"
	}
	return "unknown"
}

func ActionFromString(s string) ActionType {
	switch s {
	case "shutdown":
		return ActionShutdown
	case "restart":
		return ActionRestart
	case "lock":
		return ActionLock
	case "command":
		return ActionCommand
	}
	return ActionCommand
}

type RepeatType int

const (
	RepeatOnce RepeatType = iota
	RepeatDaily
	RepeatWeekly
	RepeatHourly
)

type MissedPolicy int

const (
	MissedSkip    MissedPolicy = iota // 错过就跳过
	MissedExecute                     // 错过立即执行
)

type Task struct {
	ID        string     `json:"id"`
	CreatedAt time.Time  `json:"created_at"`

	Action  ActionType `json:"action"`
	Command string     `json:"command,omitempty"` // ActionCommand 时

	TargetTime time.Time `json:"target_time"` // 首次执行时间

	Repeat   RepeatType `json:"repeat"`
	MaxCount int        `json:"max_count"`   // 1=单次（隐藏重复选项）, 0=无限
	Executed int        `json:"executed"`    // 已执行次数

	// 通知
	NotifyMin int  `json:"notify_min"` // 提前分钟，0=不通知
	Important bool `json:"important"`  // 重要通知（专注模式也显示）

	// 持久化策略
	Persistent bool `json:"persistent"` // true=持久化保留, false=重启销毁

	// 错过策略
	MissedPolicy MissedPolicy `json:"missed_policy"`

	Enabled bool `json:"enabled"`
}

// NextRun 计算下次执行时间
// 单次任务：返回 TargetTime（如果已到且错过，按 MissedPolicy 决定）
// 重复任务：基于已执行次数递推
func (t *Task) NextRun() time.Time {
	if t.Repeat == RepeatOnce || t.MaxCount == 1 {
		return t.TargetTime
	}

	now := time.Now()
	base := t.TargetTime

	switch t.Repeat {
	case RepeatDaily:
		// 每天同一时间（HH:MM）
		for i := t.Executed; ; i++ {
			next := time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), 0, 0, base.Location())
			next = next.AddDate(0, 0, i)
			if next.After(now) || next.Equal(now) {
				return next
			}
			if t.MaxCount > 0 && i >= t.MaxCount {
				break
			}
		}
	case RepeatWeekly:
		for i := t.Executed; ; i++ {
			next := time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), 0, 0, base.Location())
			next = next.AddDate(0, 0, i*7)
			if next.After(now) || next.Equal(now) {
				return next
			}
			if t.MaxCount > 0 && i >= t.MaxCount {
				break
			}
		}
	case RepeatHourly:
		// 每小时：每次迭代加 1 小时，保留时间模式
		for i := t.Executed; ; i++ {
			next := time.Date(base.Year(), base.Month(), base.Day(), base.Hour(), base.Minute(), 0, 0, base.Location())
			next = next.Add(time.Duration(i) * time.Hour)
			if next.After(now) || next.Equal(now) {
				return next
			}
			if t.MaxCount > 0 && i >= t.MaxCount {
				break
			}
		}
	}

	// 超出最大次数或无法计算，返回 zero
	return time.Time{}
}

// ShouldNotify 检查现在是否应该发通知
func (t *Task) ShouldNotify(notifiedKey map[string]bool) bool {
	if t.NotifyMin <= 0 {
		return false
	}
	next := t.NextRun()
	if next.IsZero() {
		return false
	}
	notifyAt := next.Add(-time.Duration(t.NotifyMin) * time.Minute)
	now := time.Now()

	if now.After(notifyAt) && now.Before(next) {
		key := t.ID + "-notify-" + next.Format("1504")
		if !notifiedKey[key] {
			notifiedKey[key] = true
			return true
		}
	}
	return false
}
