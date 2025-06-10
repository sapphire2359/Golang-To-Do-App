package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"
)

// status constants
const (
	StatusNotStarted = "not started"
	StatusStarted    = "started"
	StatusComplete   = "complete"
)

// for validation of status
var validStatuses = map[string]bool{
	StatusNotStarted: true,
	StatusStarted:    true,
	StatusComplete:   true,
}

// TodoItem struct type
type TodoItem struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

// json file where all items are going to be stored
const filename = "data/todo.json"

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
		// error handler when user does not input "desc" or "status"
		if *desc == "" || *status == "" {
			fmt.Println("Use -desc and -status with -add")
			return
		}
		addTodoItem(*desc, strings.ToLower(*status))

	case *list:
		listTodoItems()

	case *update > 0:
		// error handler when user does not input "desc" and "status"
		if *desc == "" && *status == "" {
			fmt.Println(" Usage: -update=ID -desc=\"new task\" -status=\"started\"")
			return
		}
		updateTodoItem(*update, *desc, strings.ToLower(*status))

	case *delete > 0:
		deleteTodoItem(*delete)

	default:
		fmt.Println("Usage:")
		fmt.Println("  -add -desc=\"...\" -status=\"not started|started|completed\"")
		fmt.Println("  -list")
		fmt.Println("  -update=Id -desc=\"...\" -status=\"...\"")
		fmt.Println("  -delete=Id")
	}
}

/*
[1] Checks if there is error when reading the json file
[2] Checks if the the todo item list exists
[3] Unmarshalls the json data from the file and loads it to memory
*/
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

/*
[1] Checks if there is error when marshalling the todo items passed from memory
and returns error if there is problem
[2] Writes todo items to the specified file from memory and throws error if there is problem
*/
func saveTodoItems(items []TodoItem) error {
	data, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

/*
[1] Finds out the index of the maximum or last todo item
[2] Returns the index of the item after the maximum todo item
*/
func getNextId(items []TodoItem) int {
	maxId := 0
	for _, item := range items {
		if item.Id > maxId {
			maxId = item.Id
		}
	}
	return maxId + 1
}

/*
[1] Checks if the "status" of the new todo item is valid
[2] assigns all todo items to "items" variable
[3] new todo item is created using the passed "desc" and "status" and new Id
*/
func addTodoItem(desc, status string) {
	// validation to ensure only "not started", "started", "complete" is input
	if !validStatuses[status] {
		fmt.Println("Invalid status. Use: not started, started, or complete.")
		return
	}
	items, _ := loadTodoItems()
	//items = append(items, TodoItem{Description: desc, Status: status})
	newItem := TodoItem{
		Id:          getNextId(items),
		Description: desc,
		Status:      status,
	}
	items = append(items, newItem)

	saveTodoItems(items)
	fmt.Printf("Added item [%d]: %s (%s)\n", newItem.Id, newItem.Description, newItem.Status)
}

// List all to-do items
func listTodoItems() {
	items, err := loadTodoItems()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if len(items) == 0 {
		fmt.Println("No todo items found. List is empty...")
		return
	}
	for _, item := range items {
		fmt.Printf("[%d] %s [%s]\n", item.Id, item.Description, item.Status)
	}
}

// Update an existing to-do item
func updateTodoItem(index int, desc, status string) {
	items, _ := loadTodoItems()
	updated := false
	for i, item := range items { // look for item in the item list
		if item.Id == index {
			if desc != "" {
				items[i].Description = desc
			}
			if status != "" {
				//status = strings.ToLower(status)
				if !validStatuses[status] {
					fmt.Println("Invalid status. Use: not started, started, or complete.")
					return
				}
				items[i].Status = status
			}
			updated = true
			break
		}
	}
	if !updated {
		fmt.Println("Item not found.")
		return
	}
	saveTodoItems(items)
	fmt.Println("Item updated.")
}

// Delete an item
func deleteTodoItem(id int) {
	items, _ := loadTodoItems()
	newItems := []TodoItem{}
	deleted := false
	for _, item := range items {
		if item.Id == id {
			deleted = true
			continue
		}
		newItems = append(newItems, item)
	}
	if !deleted {
		fmt.Println("Item not found.")
		return
	}
	saveTodoItems(newItems)
	fmt.Println("Item deleted.")
}
