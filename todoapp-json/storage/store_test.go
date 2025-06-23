package storage

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"
	"todoapp-json/models"
)

const testDataFile = "data/data.json" // Use the actual data file name your StoreLoop writes to

// **********************************
// Parallel test
// **********************************
func TestStoreLoopParallel(t *testing.T) {
	if err := os.Remove(testDataFile); err != nil && !os.IsNotExist(err) {
		t.Fatalf("Failed to remove test file: %v", err)
	}

	StoreChan = make(chan StoreCommand, 100)
	go StartStoreLoop()

	const count = 20

	t.Run("AddTodos", func(t *testing.T) {
		for i := 0; i < count; i++ {
			i := i
			t.Run(fmt.Sprintf("Add-%d", i), func(t *testing.T) {
				t.Parallel()

				respCh := make(chan StoreResponse, 1)
				StoreChan <- StoreCommand{
					Type: CmdAdd,
					Item: models.TodoItem{
						Description: fmt.Sprintf("Task %d", i),
						Status:      "not started",
					},
					ResponseCh: respCh,
				}

				select {
				case resp := <-respCh:
					if resp.Err != nil {
						t.Errorf("Add failed: %v", resp.Err)
					}
				case <-time.After(2 * time.Second):
					t.Fatal("Timeout waiting for add")
				}
			})
		}
	})
	// This runs AFTER all "AddTodos" children finish
	t.Run("ListTodos", func(t *testing.T) {
		respCh := make(chan StoreResponse, 1)
		StoreChan <- StoreCommand{
			Type:       CmdList,
			ResponseCh: respCh,
		}

		select {
		case resp := <-respCh:
			if resp.Err != nil {
				t.Fatalf("List failed: %v", resp.Err)
			}
			if len(resp.TodoItems) != count {
				t.Errorf("Expected %d todos, got %d", count, len(resp.TodoItems))
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for List response")
		}
	})

	t.Run("UpdateTodos", func(t *testing.T) {
		for i := 0; i < count; i++ {
			i := i
			t.Run(fmt.Sprintf("Update-%d", i), func(t *testing.T) {
				t.Parallel()

				respCh := make(chan StoreResponse, 1)
				StoreChan <- StoreCommand{
					Type: CmdUpdate,
					Item: models.TodoItem{
						Id:          i + 1, // assuming IDs are sequential starting at 1
						Description: fmt.Sprintf("Updated Task %d", i),
						Status:      "in progress",
					},
					ResponseCh: respCh,
				}

				select {
				case resp := <-respCh:
					if resp.Err != nil {
						t.Errorf("Update failed: %v", resp.Err)
					}
				case <-time.After(2 * time.Second):
					t.Fatal("Timeout waiting for update")
				}
			})
		}
	})

	t.Run("DeleteTodos", func(t *testing.T) {
		for i := 0; i < count; i++ {
			i := i
			t.Run(fmt.Sprintf("Delete-%d", i), func(t *testing.T) {
				t.Parallel()

				respCh := make(chan StoreResponse, 1)
				StoreChan <- StoreCommand{
					Type: CmdDelete,
					Id:   i + 1,
					// Item:       models.TodoItem{Id: i + 1},
					ResponseCh: respCh,
				}

				select {
				case resp := <-respCh:
					if resp.Err != nil {
						t.Errorf("Delete failed: %v", resp.Err)
					}
				case <-time.After(2 * time.Second):
					t.Fatal("Timeout waiting for delete")
				}
			})
		}
	})

	t.Run("ListShouldBeEmpty", func(t *testing.T) {
		respCh := make(chan StoreResponse, 1)
		StoreChan <- StoreCommand{Type: CmdList, ResponseCh: respCh}
		select {
		case resp := <-respCh:
			if resp.Err != nil {
				t.Fatalf("List failed: %v", resp.Err)
			}
			if len(resp.TodoItems) != 0 {
				t.Errorf("Expected 0 todos, got %d", len(resp.TodoItems))
			}
		case <-time.After(2 * time.Second):
			t.Fatal("Timeout waiting for List after delete")
		}
	})
}

// **********************************
// Unit test
// **********************************
func send(t *testing.T, cmd StoreCommand) StoreResponse {
	t.Helper()
	respCh := make(chan StoreResponse, 1)
	cmd.ResponseCh = respCh
	StoreChan <- cmd

	select {
	case resp := <-respCh:
		return resp
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for response")
		return StoreResponse{Err: context.DeadlineExceeded}
	}
}

func setupTestStoreLoop(t *testing.T) {
	t.Helper()
	SetDataFilePath(testDataFile)

	// Clean up any existing data
	_ = os.Remove(testDataFile)

	StoreChan = make(chan StoreCommand, 100)
	go StartStoreLoop()
}

func TestStartStoreLoop_Behaviors(t *testing.T) {
	setupTestStoreLoop(t)

	// Add item
	item := models.TodoItem{
		Description: "Test task",
		Status:      models.NotStarted,
	}
	addResp := send(t, StoreCommand{
		Type: CmdAdd,
		Item: item,
	})
	if addResp.Err != nil {
		t.Fatalf("Add failed: %v", addResp.Err)
	}

	// List items
	listResp := send(t, StoreCommand{Type: CmdList})
	if listResp.Err != nil {
		t.Fatalf("List failed: %v", listResp.Err)
	}
	if len(listResp.TodoItems) != 1 {
		t.Fatalf("Expected 1 item, got %d", len(listResp.TodoItems))
	}
	if listResp.TodoItems[0].Description != "Test task" {
		t.Errorf("Unexpected item: %+v", listResp.TodoItems[0])
	}

	// Update item
	updatedItem := listResp.TodoItems[0]
	updatedItem.Description = "Updated task"
	updatedItem.Status = models.Completed

	updateResp := send(t, StoreCommand{
		Type: CmdUpdate,
		Item: updatedItem,
	})
	if updateResp.Err != nil {
		t.Fatalf("Update failed: %v", updateResp.Err)
	}

	// Delete item
	delResp := send(t, StoreCommand{
		Type: CmdDelete,
		Id:   updatedItem.Id,
	})
	if delResp.Err != nil {
		t.Fatalf("Delete failed: %v", delResp.Err)
	}

	// Verify empty
	finalList := send(t, StoreCommand{Type: CmdList})
	if len(finalList.TodoItems) != 0 {
		t.Fatalf("Expected 0 items, got %d", len(finalList.TodoItems))
	}
}

// func TestStartStoreLoop_UpdateNonexistent(t *testing.T) {
// 	setupTestStoreLoop(t)

// 	nonexistent := models.TodoItem{Id: 999, Description: "X", Status: models.NotStarted}
// 	resp := send(t, StoreCommand{
// 		Type: CmdUpdate,
// 		Item: nonexistent,
// 	})
// 	if resp.Err == nil || !errors.Is(resp.Err, errors.New("update : todo item not found")) {
// 		t.Fatalf("Expected not found error, got %v", resp.Err)
// 	}
// }

// func TestStartStoreLoop_DeleteNonexistent(t *testing.T) {
// 	setupTestStoreLoop(t)

// 	resp := send(t, StoreCommand{
// 		Type: CmdDelete,
// 		Id:   888,
// 	})
// 	if resp.Err == nil || !errors.Is(resp.Err, errors.New("delete : todo item not found")) {
// 		t.Fatalf("Expected delete not found error, got %v", resp.Err)
// 	}
// }
