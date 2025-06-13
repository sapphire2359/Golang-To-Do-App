package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// file where items are stored
const fileName = "data/todo-items.txt"

// displays all todo item option
func displayTodoOptions() {
	fmt.Println("\n********************To Do Options*******************")
	fmt.Println("1. Add Item")
	fmt.Println("2. Display Items")
	fmt.Println("3. Update Item Status")
	fmt.Println("4. Delete Item")
	fmt.Println("5. Exit")
	fmt.Println("******************************************************")
}

// adds todo item
func addToDoItem() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Add To-Do Item")

	// Get item description
	fmt.Print("Enter item description: ")
	description, _ := reader.ReadString('\n')
	description = strings.TrimSpace(description)

	// Get item status
	fmt.Print("Enter item status [not started | started | completed]: ")
	status, _ := reader.ReadString('\n')
	status = strings.TrimSpace(strings.ToLower(status))

	// status validation
	validStatuses := map[string]bool{"not started": true, "started": true, "completed": true}
	if !validStatuses[status] {
		fmt.Println("Invalid status entered.")
		return
	}

	// Open file in append mode, create if not exists
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write formatted line
	//line := fmt.Sprintf("Item Description: %s, Item Status: %s\n", description, status)
	line := fmt.Sprintf("%s, %s\n", description, status)
	if _, err := file.WriteString(line); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("To-Do item added successfully.")
}

// displays all todo items
func displayAllTodoItems() {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	fmt.Println("\nTo Do Items List:")
	scanner := bufio.NewScanner(file)
	index := 1
	for scanner.Scan() {
		parts := strings.Split(scanner.Text(), ",")
		if len(parts) == 2 {
			//fmt.Printf("%d. Description: %s | Status: %s\n", index, parts[0], parts[1])
			fmt.Printf("%d. %s [%s ]\n", index, parts[0], parts[1])
			index++
		}
	}
}

// update todo item status or description
func updateTodoItemStatus() {
	displayAllTodoItems()

	fmt.Print("Enter item number to update: ")
	var num int
	fmt.Scanln(&num)

	file, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("File read error:", err)
		return
	}

	lines := strings.Split(string(file), "\n")
	if num < 1 || num > len(lines)-1 {
		fmt.Println("Invalid item number.")
		return
	}

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter new description: ")
	newDesc, _ := reader.ReadString('\n')
	newDesc = strings.TrimSpace(newDesc)

	fmt.Print("Enter new status: ")
	newStatus, _ := reader.ReadString('\n')
	newStatus = strings.TrimSpace(newStatus)

	// status validation
	validStatuses := map[string]bool{"not started": true, "started": true, "completed": true}
	if !validStatuses[newStatus] {
		fmt.Println("Invalid status entered.")
		return
	}

	lines[num-1] = fmt.Sprintf("%s, %s", newDesc, newStatus)
	//lines[num-1] = fmt.Sprintf("Description: %s, Status: %s\n", newDesc, newStatus)

	err = os.WriteFile(fileName, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}
	fmt.Println("Item updated.")
}

// delete todo item status or description
func deleteTodoItem() {
	displayAllTodoItems()

	fmt.Print("Enter item number to delete: ")
	var num int
	fmt.Scanln(&num)

	// Read all lines
	content, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Println("Could not read file:", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	// Validate item number
	if num < 1 || num > len(lines)-1 {
		fmt.Println("Invalid item number.")
		return
	}

	// Remove the specified item
	lines = append(lines[:num-1], lines[num:]...)

	// Join and write updated lines back to file
	newContent := strings.Join(lines, "\n")
	err = os.WriteFile(fileName, []byte(newContent), 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}
	fmt.Println("Item deleted successfully.")
}

func main() {
	for {
		var todoOption int

		displayTodoOptions()
		fmt.Println("Select To Do option [1 to 5]: ")
		fmt.Scanln(&todoOption)

		//i := 2
		switch todoOption {
		case 1:
			addToDoItem()

		case 2:
			displayAllTodoItems()

		case 3:
			updateTodoItemStatus()

		case 4:
			deleteTodoItem()

		case 5:
			fmt.Println("Exiting app...")
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}
