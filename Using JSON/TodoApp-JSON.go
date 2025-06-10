package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

type TodoItem struct {
	Description string `json:"description"`
	Status      string `json:"status"`
}

// json file where all items are going to be stored
const filename = "todo.json"

// for validation of status
var validStatuses = map[string]bool{
	"not started": true,
	"started":     true,
	"complete":    true,
}

// Load existing items from JSON file
func loadTodoItems() ([]TodoItem, error) {
	var items []TodoItem
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []TodoItem{}, nil // return empty list if file doesn't exist
		}
		return nil, err
	}
	err = json.Unmarshal(data, &items)
	return items, err
}

// Save items to JSON file
func saveTodoItems(items []TodoItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// Add a new to-do item
func addTodoItem(desc, status string) {
	if !validStatuses[status] {
		fmt.Println("Invalid status. Use: not started, started, or complete.")
		return
	}
	items, _ := loadTodoItems()
	items = append(items, TodoItem{Description: desc, Status: status})

	saveTodoItems(items)
	fmt.Println("Item added.")
}

// List all to-do items
func listTodoItems() {
	items, err := loadTodoItems()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(items) == 0 {
		fmt.Println("No to-do items found.")
		return
	}
	for i, item := range items {
		fmt.Printf("%d. %s [%s]\n", i+1, item.Description, item.Status)
	}
}

// Update an existing to-do item
func updateTodoItem(index int, desc, status string) {
	items, _ := loadTodoItems()
	if index < 1 || index > len(items) {
		fmt.Println("Invalid item number.")
		return
	}
	if desc != "" {
		items[index-1].Description = desc
	}
	if status != "" {
		if !validStatuses[status] {
			fmt.Println("Invalid status. Use: not started, started, or complete.")
			return
		}
		items[index-1].Status = status
	}
	saveTodoItems(items)
	fmt.Println("Item updated.")
}

// Delete an item
func deleteTodoItem(index int) {
	items, _ := loadTodoItems()
	if index < 1 || index > len(items) {
		fmt.Println("Invalid item number.")
		return
	}
	items = append(items[:index-1], items[index:]...)
	saveTodoItems(items)
	fmt.Println("Item deleted.")
}

func main() {
	add := flag.Bool("add", false, "Add a new todo item")
	list := flag.Bool("list", false, "List all todo items")
	update := flag.Int("update", 0, "Update item by number")
	delete := flag.Int("delete", 0, "Delete item by number")

	desc := flag.String("desc", "", "Description of the todo item")
	status := flag.String("status", "", "Todo status (not started | started | complete)")

	flag.Parse()

	switch {
	case *add:
		if *desc == "" || *status == "" {
			fmt.Println("Use -desc and -status with -add")
			return
		}
		addTodoItem(*desc, strings.ToLower(*status))

	case *list:
		listTodoItems()

	case *update > 0:
		updateTodoItem(*update, *desc, strings.ToLower(*status))

	case *delete > 0:
		deleteTodoItem(*delete)

	default:
		fmt.Println("Usage:")
		fmt.Println("  -add -desc=\"...\" -status=\"not started|started|completed\"")
		fmt.Println("  -list")
		fmt.Println("  -update=N -desc=\"...\" -status=\"...\"")
		fmt.Println("  -delete=N")
	}
}
