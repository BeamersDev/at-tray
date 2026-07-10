package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
)

// Storage 管理任务的持久化
type Storage struct {
	mu     sync.Mutex
	file   string
	tasks  []*Task
}

func NewStorage() (*Storage, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		appData = filepath.Join(os.Getenv("HOME"), ".at-tray")
	}
	dir := filepath.Join(appData, "at-tray")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}
	s := &Storage{
		file:  filepath.Join(dir, "tasks.json"),
		tasks: []*Task{},
	}
	s.load()
	return s, nil
}

func (s *Storage) load() {
	data, err := os.ReadFile(s.file)
	if err != nil {
		return // 文件不存在或读不了就空列表
	}
	_ = json.Unmarshal(data, &s.tasks)
	if s.tasks == nil {
		s.tasks = []*Task{}
	}
}

func (s *Storage) save() {
	data, _ := json.MarshalIndent(s.tasks, "", "  ")
	_ = os.WriteFile(s.file, data, 0644)
}

func (s *Storage) All() []*Task {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.tasks
}

func (s *Storage) Add(t *Task) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.tasks = append(s.tasks, t)
	s.save()
}

func (s *Storage) Delete(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for i, t := range s.tasks {
		if t.ID == id {
			s.tasks = append(s.tasks[:i], s.tasks[i+1:]...)
			s.save()
			return true
		}
	}
	return false
}

func (s *Storage) Update(id string, fn func(t *Task)) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, t := range s.tasks {
		if t.ID == id {
			fn(t)
			s.save()
			return true
		}
	}
	return false
}
