package main

import (
	"os"
	"testing"
)

func TestLoadTasks(t *testing.T) {
	t.Run("File exists with data", func(t *testing.T) {
		store := TaskStore{Tasks: []Task{{ID: 1, Description: "Test task", Done: false}}}
		if err := SaveTasks(store); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		loaded, err := LoadTasks()
		if err != nil {
			t.Errorf("Failed to load tasks: %v", err)
		}
		if len(loaded.Tasks) == 0 {
			t.Errorf("Expected tasks, got empty store")
		}
		if loaded.Tasks[0].Description != "Test task" {
			t.Errorf("Expected 'Test task', got %v", loaded.Tasks[0].Description)
		}
	})

	t.Run("File exists but empty", func(t *testing.T) {
		if err := os.WriteFile(taskFile, []byte{}, 0644); err != nil {
			t.Fatalf("Setup failed: %v", err)
		}

		loaded, err := LoadTasks()
		if err != nil {
			t.Errorf("Failed to load tasks: %v", err)
		}
		if len(loaded.Tasks) != 0 {
			t.Errorf("Expected empty store, got %v", loaded.Tasks)
		}
	})

	t.Run("File does not exist", func(t *testing.T) {
		os.Remove(taskFile)

		loaded, err := LoadTasks()
		if err != nil {
			t.Errorf("Failed to load tasks: %v", err)
		}
		if len(loaded.Tasks) != 0 {
			t.Errorf("Expected empty store, got %v", loaded.Tasks)
		}
	})
}
