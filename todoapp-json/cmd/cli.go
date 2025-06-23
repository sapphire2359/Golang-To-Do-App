package cmd

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"todoapp-json/logic"
	"todoapp-json/models"
)

func RunCLI(ctx context.Context) {
	description := flag.String("desc", "", "Description for the todo item")
	status := flag.String("status", "not started", "Status [not started|started|completed]")
	action := flag.String("action", "list", "Action to perform [add|list|update|delete]")
	id := flag.Int("id", 0, "ID of the item to update/delete")
	flag.Parse()

	switch *action {
	case "add":
		addTodoItem(ctx, *description, models.Status(*status))

	case "list":
		listTodoItems(ctx)

	case "update":
		updateTodoItem(ctx, *id, *description, models.Status(*status))

	case "delete":
		deleteTodoItem(ctx, *id)

	default:
		fmt.Println("How to use Todo App CLI :")
		fmt.Println("  -action=add -desc=\"...\" -status=\"not started|started|completed\"")
		fmt.Println("  -action=list")
		fmt.Println("  -action=update -Id=\"...\" -desc=\"...\" -status=\"...\"")
		fmt.Println("  -action=delete= -Id=\"...\"")
	}
}

// addTodoItem
func addTodoItem(ctx context.Context, description string, status models.Status) {
	if strings.TrimSpace(description) == "" {
		fmt.Println("Error: -action=add -desc=\"...\" -status=\"not started|started|completed\"")
		return
	}
	if !models.IsValidStatus(models.Status(status)) {
		fmt.Println("Error: invalid status. input only [not started|started|completed]")
		return
	}
	err := logic.AddTodoItem(ctx, description, models.Status(status))
	if err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("Todo item added...")
	}
}

// listTodoItems
func listTodoItems(ctx context.Context) {
	todoItems, err := logic.ListTodoItems(ctx)
	if err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		for _, item := range todoItems {
			fmt.Printf("[%d] %s [%s]\n", item.Id, item.Description, item.Status)
		}
	}
}

// update todo item
func updateTodoItem(ctx context.Context, id int, description string, status models.Status) {
	if id == 0 || strings.TrimSpace(description) == "" {
		fmt.Println("Error: -action=update -id=\"...\" -desc=\"...\" -status=\"...\"")
		return
	}
	if !models.IsValidStatus(models.Status(status)) {
		fmt.Println("Error: invalid status. input only [not started|started|completed]")
		return
	}

	err := logic.UpdateTodoItem(ctx, id, description, models.Status(status))
	if err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("Todo item updated...")
	}
}

// delete todo item
func deleteTodoItem(ctx context.Context, id int) {
	if id == 0 {
		fmt.Println("Error: -action=delete= -id=\"...\"")
		return
	}
	err := logic.DeleteTodoItem(ctx, id)
	if err != nil {
		fmt.Println("Error:", err.Error())
	} else {
		fmt.Println("Todo item deleted...")
	}
}
