package logic_test

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"
	"todoapp-json/logic"
	"todoapp-json/models"
	"todoapp-json/storage"
)

const testDataFile = "data/data.json"

// helper to send a command and wait for a response
func sendCommand(cmd storage.StoreCommand) (resp storage.StoreResponse) {
	respCh := make(chan storage.StoreResponse, 1)
	cmd.ResponseCh = respCh
	storage.StoreChan <- cmd

	select {
	case resp = <-respCh:
		return
	case <-time.After(2 * time.Second):
		return storage.StoreResponse{Err: context.DeadlineExceeded}
	}
}

// ****************
// Parallel test
// ***************
func TestStoreIntegrationParallel(t *testing.T) {
	// testDataFile := filepath.Join(t.TempDir(), "todos_test.json")
	storage.SetDataFilePath(testDataFile)

	// Clean up any existing test file
	if err := os.Remove(testDataFile); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove test file: %v", err)
	}

	// Start store loop fresh
	storage.StoreChan = make(chan storage.StoreCommand, 100)
	go storage.StartStoreLoop()

	t.Run("ConcurrentAccess", func(t *testing.T) {
		const count = 100

		t.Run("AddTodosParallel", func(t *testing.T) {
			for i := 0; i < count; i++ {
				i := i
				t.Run(fmt.Sprintf("Add-%d", i), func(t *testing.T) {
					t.Parallel()
					err := logic.AddTodoItem(context.Background(), fmt.Sprintf("Task %d", i), models.NotStarted)
					if err != nil {
						t.Errorf("Failed to add task %d: %v", i, err)
					}
				})
			}
		})

		// Wait for all adds to complete
		t.Run("ListToVerifyAdd", func(t *testing.T) {
			resp := sendCommand(storage.StoreCommand{Type: storage.CmdList})
			if resp.Err != nil {
				t.Fatalf("List after add failed: %v", resp.Err)
			}
			if len(resp.TodoItems) != count {
				t.Fatalf("Expected %d todos, got %d", count, len(resp.TodoItems))
			}
		})

		// Update all in parallel
		t.Run("UpdateTodosParallel", func(t *testing.T) {
			resp := sendCommand(storage.StoreCommand{Type: storage.CmdList})
			if resp.Err != nil {
				t.Fatalf("List before update failed: %v", resp.Err)
			}
			for _, item := range resp.TodoItems {
				item := item
				t.Run(fmt.Sprintf("Update-%d", item.Id), func(t *testing.T) {
					t.Parallel()
					item.Description += " (updated)"
					item.Status = models.Completed
					resp := sendCommand(storage.StoreCommand{
						Type: storage.CmdUpdate,
						Item: item,
					})
					if resp.Err != nil {
						t.Errorf("Update failed for ID %d: %v", item.Id, resp.Err)
					}
				})
			}
		})

		// Delete all in parallel
		t.Run("DeleteTodosParallel", func(t *testing.T) {
			resp := sendCommand(storage.StoreCommand{Type: storage.CmdList})
			if resp.Err != nil {
				t.Fatalf("List before delete failed: %v", resp.Err)
			}
			for _, item := range resp.TodoItems {
				item := item
				t.Run(fmt.Sprintf("Delete-%d", item.Id), func(t *testing.T) {
					t.Parallel()
					resp := sendCommand(storage.StoreCommand{
						Type: storage.CmdDelete,
						Id:   item.Id,
					})
					if resp.Err != nil {
						t.Errorf("Delete failed for ID %d: %v", item.Id, resp.Err)
					}
				})
			}
		})

		// Verify all are gone
		t.Run("ListAfterAllDeletes", func(t *testing.T) {
			resp := sendCommand(storage.StoreCommand{Type: storage.CmdList})
			if resp.Err != nil {
				t.Fatalf("List failed after deletes: %v", resp.Err)
			}
			if len(resp.TodoItems) != 0 {
				t.Errorf("Expected 0 todos after deletes, got %d", len(resp.TodoItems))
			}
		})
	})
}

// ***************
// Unit test
// ***************
func TestStoreIntegrationUnit(t *testing.T) {
	storage.SetDataFilePath(testDataFile) // ensure your storage package supports setting file path

	// Explicitly remove file before starting
	if err := os.Remove(testDataFile); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove test file: %v", err)
	}
	// Start real store loop
	storage.StoreChan = make(chan storage.StoreCommand, 100)
	go storage.StartStoreLoop()

	totalTodoItems := 5

	// 1. Add Todo
	t.Run("AddTodo", func(t *testing.T) {
		for i := 0; i < totalTodoItems; i++ {
			err := logic.AddTodoItem(context.Background(), "Wash dishes "+strconv.Itoa(i), models.NotStarted)
			if err != nil {
				t.Fatalf("AddTodoItem failed: %v", err)
			}
		}
	})

	// 2. List Todos
	var list []models.TodoItem
	t.Run("ListTodos", func(t *testing.T) {
		resp := sendCommand(storage.StoreCommand{Type: storage.CmdList})
		if resp.Err != nil {
			t.Fatalf("List failed: %v", resp.Err)
		}
		if len(resp.TodoItems) != totalTodoItems {
			t.Fatalf("Expected 1 item, got %d", len(resp.TodoItems))
		}
		list = resp.TodoItems
	})

	// 3. Update Todo
	t.Run("UpdateTodo", func(t *testing.T) {
		item := list[0]
		item.Description = "Wash dishes thoroughly"
		item.Status = models.Completed
		resp := sendCommand(storage.StoreCommand{
			Type: storage.CmdUpdate,
			Item: item,
		})
		if resp.Err != nil {
			t.Fatalf("Update failed: %v", resp.Err)
		}
	})

	// 4. Delete Todo
	t.Run("DeleteTodo", func(t *testing.T) {
		item := list[0]

		resp := sendCommand(storage.StoreCommand{
			Type: storage.CmdDelete,
			Id:   item.Id,
		})
		if resp.Err != nil {
			t.Fatalf("Delete failed: %v", resp.Err)
		}
	})

	// 5. Verify List is Empty
	t.Run("ListAfterDelete", func(t *testing.T) {
		resp := sendCommand(storage.StoreCommand{Type: storage.CmdList})
		if resp.Err != nil {
			t.Fatalf("Final list failed: %v", resp.Err)
		}
		if len(resp.TodoItems) != totalTodoItems-1 {
			t.Errorf("Expected 0 items after delete, got %d", len(resp.TodoItems))
		}
	})
}
