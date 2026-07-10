package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gogpu/systray"
)

var (
	storage     *Storage
	srv         *http.Server
	port        int
	tray        *systray.SystemTray
	icon        []byte
	iconDark    []byte
	notifiedKey = make(map[string]bool)
)

func main() {
	var err error
	storage, err = NewStorage()
	if err != nil {
		log.Fatalf("初始化存储失败: %v", err)
	}

	icon = genIcon(22, 200, 200, 200)
	iconDark = genIcon(22, 80, 200, 255)

	// 启动 HTTP server
	startServer()

	// 启动任务调度
	go schedulerLoop()

	// 启动托盘
	tray = systray.New()
	tray.SetIcon(icon).SetDarkModeIcon(iconDark).
		SetTooltip("at-tray — 定时任务管理器").
		SetMenu(buildMenu()).
		Show()
	tray.Run()

	// 程序退出 —— 清理非持久化任务
	cleanupNonPersistent()
}

func cleanupNonPersistent() {
	for _, t := range storage.All() {
		if !t.Persistent {
			storage.Delete(t.ID)
		}
	}
}

func buildMenu() *systray.Menu {
	m := systray.NewMenu()
	m.Add("📋 打开管理页面", func() {
		url := fmt.Sprintf("http://127.0.0.1:%d", port)
		exec.Command("cmd", "/c", "start", url).Start()
	})
	m.AddSeparator()
	m.Add("⏹ 退出", func() {
		if srv != nil {
			srv.Shutdown(context.Background())
		}
		cleanupNonPersistent()
		tray.Remove()
		os.Exit(0)
	})
	return m
}

func startServer() {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatalf("监听端口失败: %v", err)
	}
	port = listener.Addr().(*net.TCPAddr).Port

	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.HandleFunc("/api/tasks", handleTasks)
	mux.HandleFunc("/api/tasks/", handleTaskByID)
	mux.HandleFunc("/api/execute", handleExecute)
	mux.HandleFunc("/api/notify", handleNotify)
	mux.HandleFunc("/api/status", handleStatus)

	srv = &http.Server{Handler: mux}
	go srv.Serve(listener)

	log.Printf("管理页面: http://127.0.0.1:%d", port)
}

// ── HTTP Handlers ──

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(indexHTML))
}

func handleTasks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "GET":
		json.NewEncoder(w).Encode(storage.All())

	case "POST":
		var t Task
		if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
			http.Error(w, `{"error":"invalid json"}`, 400)
			return
		}
		t.ID = fmt.Sprintf("%d", time.Now().UnixNano())
		t.CreatedAt = time.Now()
		t.Executed = 0
		t.MaxCount = 0 // 0 = unlimited
		t.Enabled = true
		storage.Add(&t)
		w.WriteHeader(201)
		json.NewEncoder(w).Encode(t)

	default:
		http.Error(w, "", 405)
	}
}

func handleTaskByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id := strings.TrimPrefix(r.URL.Path, "/api/tasks/")

	// 如果 id 为空（匹配 /api/tasks/ 但没有具体 ID），返回 400
	if id == "" {
		http.Error(w, `{"error":"missing id"}`, 400)
		return
	}

	switch r.Method {
	case "GET":
		for _, t := range storage.All() {
			if t.ID == id {
				json.NewEncoder(w).Encode(t)
				return
			}
		}
		http.Error(w, `{"error":"not found"}`, 404)

	case "DELETE":
		if storage.Delete(id) {
			json.NewEncoder(w).Encode(map[string]bool{"ok": true})
		} else {
			http.Error(w, `{"error":"not found"}`, 404)
		}

	case "PATCH":
		var body struct {
			Enabled      *bool        `json:"enabled,omitempty"`
			Action       *ActionType  `json:"action,omitempty"`
			Command      *string      `json:"command,omitempty"`
			TargetTime   *time.Time   `json:"target_time,omitempty"`
			Repeat       *RepeatType  `json:"repeat,omitempty"`
			MaxCount     *int         `json:"max_count,omitempty"`
			Executed     *int         `json:"executed,omitempty"`
			NotifyMin    *int         `json:"notify_min,omitempty"`
			Important    *bool        `json:"important,omitempty"`
			Persistent   *bool        `json:"persistent,omitempty"`
			MissedPolicy *MissedPolicy `json:"missed_policy,omitempty"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, `{"error":"invalid json"}`, 400)
			return
		}

		updated := storage.Update(id, func(t *Task) {
			if body.Enabled != nil {
				t.Enabled = *body.Enabled
			}
			if body.Action != nil {
				t.Action = *body.Action
			}
			if body.Command != nil {
				t.Command = *body.Command
			}
			if body.TargetTime != nil {
				t.TargetTime = *body.TargetTime
			}
			if body.Repeat != nil {
				t.Repeat = *body.Repeat
			}
			if body.MaxCount != nil {
				t.MaxCount = *body.MaxCount
			}
			if body.Executed != nil {
				t.Executed = *body.Executed
			}
			if body.NotifyMin != nil {
				t.NotifyMin = *body.NotifyMin
			}
			if body.Important != nil {
				t.Important = *body.Important
			}
			if body.Persistent != nil {
				t.Persistent = *body.Persistent
			}
			if body.MissedPolicy != nil {
				t.MissedPolicy = *body.MissedPolicy
			}
		})
		if !updated {
			http.Error(w, `{"error":"not found"}`, 404)
			return
		}
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})

	default:
		http.Error(w, "", 405)
	}
}

func handleExecute(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", 405)
		return
	}
	var body struct {
		Action  ActionType `json:"action"`
		Command string     `json:"command"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid json"}`, 400)
		return
	}
	go execAction(body.Action, body.Command)
	json.NewEncoder(w).Encode(map[string]string{"status": "executed"})
}

func handleNotify(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "", 405)
		return
	}
	var body struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	}
	json.NewDecoder(r.Body).Decode(&body)
	if tray != nil {
		tray.ShowNotification(body.Title, body.Message)
	}
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func handleStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"ok":     true,
		"tasks":  len(storage.All()),
		"uptime": time.Now().Format(time.RFC3339),
	})
}

// ── Scheduler ──

func schedulerLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		for _, t := range storage.All() {
			if !t.Enabled {
				continue
			}
			next := t.NextRun()
			if next.IsZero() {
				continue
			}

			// 任务到期检查
			if !next.After(now) {
				if now.Sub(next) > 10*time.Second {
					// 真正错过了（超过10秒），按错过策略处理
					switch t.MissedPolicy {
					case MissedSkip:
						// 跳过本次，推进计数
						storage.Update(t.ID, func(up *Task) {
							up.Executed++
							if up.MaxCount > 0 && up.Executed >= up.MaxCount {
								up.Enabled = false
							}
						})
						continue
					case MissedExecute:
						// 立即执行（fall through）
					}
				}
				go execAndAdvance(t)
			}

			// 提前通知
			if t.ShouldNotify(notifiedKey) {
				if tray != nil {
					actionName := t.Action.String()
					if t.Action == ActionCommand {
						actionName = "命令: " + t.Command
					}
					msg := fmt.Sprintf("任务即将执行: %s 于 %s", actionName, next.Format("15:04"))
					tray.ShowNotification("at-tray", msg)
				}
			}
		}
	}
}

func execAndAdvance(t *Task) {
	actionName := t.Action.String()
	if t.Action == ActionCommand {
		actionName = "命令: " + t.Command
	}
	idShort := t.ID
	if len(idShort) > 8 {
		idShort = idShort[:8]
	}
	log.Printf("执行任务 %s: %s", idShort, actionName)

	execAction(t.Action, t.Command)

	// 执行后通知（如果有提前通知设置）
	if t.NotifyMin > 0 && tray != nil {
		tray.ShowNotification("at-tray", fmt.Sprintf("任务已执行: %s", actionName))
	}

	storage.Update(t.ID, func(up *Task) {
		up.Executed++
		if up.MaxCount > 0 && up.Executed >= up.MaxCount {
			up.Enabled = false
		}
	})
}

func execAction(action ActionType, command string) {
	switch action {
	case ActionShutdown:
		exec.Command("shutdown", "/s", "/t", "5").Start()
	case ActionRestart:
		exec.Command("shutdown", "/r", "/t", "5").Start()
	case ActionLock:
		exec.Command("rundll32.exe", "user32.dll,LockWorkStation").Start()
	case ActionCommand:
		if command != "" {
			parts := strings.Fields(command)
			if len(parts) > 0 {
				exec.Command(parts[0], parts[1:]...).Start()
			}
		}
	}
}

// ── Icon generation ──

func genIcon(size int, r, g, b uint8) []byte {
	return genIconPNG(size, r, g, b)
}

func genIconPNG(size int, r, g, b uint8) []byte {
	img := image.NewRGBA(image.Rect(0, 0, size, size))
	center := float64(size) / 2
	radius := float64(size)*0.45 - 1

	// 绘制实心圆
	for y := 0; y < size; y++ {
		for x := 0; x < size; x++ {
			dx := float64(x) - center
			dy := float64(y) - center
			dist := math.Sqrt(dx*dx + dy*dy)
			if dist <= radius {
				img.Set(x, y, color.RGBA{r, g, b, 255})
			} else {
				img.Set(x, y, color.RGBA{0, 0, 0, 0})
			}
		}
	}

	var buf bytes.Buffer
	_ = png.Encode(&buf, img)
	return buf.Bytes()
}
