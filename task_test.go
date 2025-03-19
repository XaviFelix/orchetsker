package main

import (
	"testing"
)

func TestLoadTasks(t *testing.T) {
	t.Run("File exists with data", func(t *testing.T) {
		loaded, err := LoadTasks()
		if err != nil {
			t.Errorf("Failed to load tasks: %v", err)
		}
		if len(loaded.Tasks) == 0 {
			t.Errorf("Expected tasks, got empty store")
		}
		entry := loaded.Tasks[1]
		if entry.Task.Description != "Test task" {
			t.Errorf("Expected 'Test task', got %v", entry.Task.Description)
		}
	})
}
