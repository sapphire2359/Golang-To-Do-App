package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"todoapp-json/models"
)

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

func RunREPL(ctx context.Context) {
	for {
		var todoOption int

		displayTodoOptions()
		fmt.Println("Select To Do option [1 to 5]: ")
		fmt.Scanln(&todoOption)

		switch todoOption {
		case 1:
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

			addTodoItem(ctx, description, models.Status(status))

		case 2:
			listTodoItems(ctx)

		case 3:
			displayTodoOptions()

			fmt.Print("Enter item number to update: ")
			var id int
			fmt.Scanln(&id)

			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter new description: ")
			newDesc, _ := reader.ReadString('\n')
			newDesc = strings.TrimSpace(newDesc)

			fmt.Print("Enter new status: ")
			newStatus, _ := reader.ReadString('\n')
			newStatus = strings.TrimSpace(newStatus)

			updateTodoItem(ctx, id, newDesc, models.Status(newStatus))

		case 4:
			displayTodoOptions()

			fmt.Print("Enter item number to delete: ")
			var id int
			fmt.Scanln(&id)

			deleteTodoItem(ctx, id)

		case 5:
			fmt.Println("Exiting app...")
			os.Exit(1)
			return

		default:
			fmt.Println("Invalid option.")
		}
	}
}
