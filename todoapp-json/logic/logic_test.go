package logic

import (
	"context"
	"os"
	"testing"
	"todoapp-json/models"
	"todoapp-json/storage"
)

const testDataFile = "data/todos.json"

// setupTest prepares a clean test environment.
func setupTest(t *testing.T) context.Context {
	t.Helper()
	err := os.WriteFile(testDataFile, []byte("[]"), 0644)
	if err != nil {
		t.Fatalf("failed to set up test data: %v", err)
	}
	return context.Background()
}

// teardownTest cleans up the test data file.
func teardownTest(t *testing.T) {
	t.Helper()
	if err := os.Remove(testDataFile); err != nil && !os.IsNotExist(err) {
		t.Logf("failed to clean up test data: %v", err)
	}
}

func TestAddTodoItem(t *testing.T) {
	ctx := setupTest(t)
	defer teardownTest(t)

	description := "Write unit tests"
	status := models.NotStarted

	err := AddTodoItem(ctx, description, status)
	if err != nil {
		t.Errorf("AddTodoItem failed: %v", err)
	}

	todoItems, err := storage.LoadTodoItems()
	if err != nil {
		t.Fatalf("LoadTodoItems failed: %v", err)
	}

	if len(todoItems) != 1 {
		t.Fatalf("expected 1 todo, got %d", len(todoItems))
	}

	if todoItems[0].Description != description || todoItems[0].Status != status {
		t.Errorf("unexpected todo item: %+v", todoItems[0])
	}
}

func TestAddTodoItem_TableDriven(t *testing.T) {
	tests := []struct {
		name        string
		description string
		status      models.Status
		wantErr     bool
	}{
		{"Valid started", "Task A", models.Started, false},
		{"Valid complete", "Task B", models.Completed, false},
		{"Valid not started", "Task C", models.NotStarted, false},
		{"Invalid status", "Bad task", models.Status("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupTest(t)
			defer teardownTest(t)

			err := AddTodoItem(ctx, tt.description, tt.status)

			if (err != nil) != tt.wantErr {
				t.Errorf("AddTodoItem() error = %v, wantErr %v", err, tt.wantErr)
			}

			todoItems, _ := storage.LoadTodoItems()
			if tt.wantErr {
				if len(todoItems) != 0 {
					t.Errorf("expected 0 items on error, got %d", len(todoItems))
				}
				return
			}
			index := len(todoItems) - 1
			// if len(todoItems) != 1 {
			// 	t.Fatalf("expected 1 item, got %d", len(todoItems))
			// }
			if todoItems[index].Description != tt.description {
				t.Errorf("expected description %q, got %q", tt.description, todoItems[index].Description)
			}
			if todoItems[index].Status != tt.status {
				t.Errorf("expected status %q, got %q", tt.status, todoItems[index].Status)
			}
		})
	}
}

func TestAddTodoItem_InvalidStatus(t *testing.T) {
	ctx := setupTest(t)
	defer teardownTest(t)

	err := AddTodoItem(ctx, "Test invalid status", "bad_status")
	if err == nil {
		t.Errorf("expected error for invalid status, got nil")
	}
}

func TestUpdateTodoItem(t *testing.T) {
	ctx := setupTest(t)
	defer teardownTest(t)

	_ = AddTodoItem(ctx, "Original task", models.NotStarted)

	err := UpdateTodoItem(ctx, 1, "Updated task", models.Completed)
	if err != nil {
		t.Errorf("UpdateTodoItem failed: %v", err)
	}

	todos, _ := storage.LoadTodoItems()
	if todos[0].Description != "Updated task" || todos[0].Status != models.Completed {
		t.Errorf("update not applied correctly: %+v", todos[0])
	}
}
func TestUpdateTodoItem_TableDriven(t *testing.T) {
	type args struct {
		id          int
		description string
		status      models.Status
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Update description", args{1, "Updated Desc", ""}, false},
		{"Update status", args{1, "", models.Completed}, false},
		{"Invalid ID", args{99, "desc", models.Started}, true},
		{"Invalid status", args{1, "", "bad_status"}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupTest(t)
			defer teardownTest(t)

			_ = AddTodoItem(ctx, "Initial", models.NotStarted)

			err := UpdateTodoItem(ctx, tt.args.id, tt.args.description, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Errorf("UpdateTodoItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDeleteTodoItem(t *testing.T) {
	ctx := setupTest(t)
	defer teardownTest(t)

	_ = AddTodoItem(ctx, "To delete", models.Started)

	err := DeleteTodoItem(ctx, 1)
	if err != nil {
		t.Errorf("DeleteTodoItem failed: %v", err)
	}

	todos, _ := storage.LoadTodoItems()
	if len(todos) != 0 {
		t.Errorf("todo not deleted, got: %+v", todos)
	}
}

func TestDeleteTodoItem_TableDriven(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{"Delete existing", 1, false},
		{"Delete non-existent", 99, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := setupTest(t)
			defer teardownTest(t)

			_ = AddTodoItem(ctx, "Task to delete", models.NotStarted)

			err := DeleteTodoItem(ctx, tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteTodoItem() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
