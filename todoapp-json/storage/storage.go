package storage

import (
	"encoding/json"
	"os"
	"todoapp-json/models"
)

const fileName = "data/todos.json"

// load todo items from json file
func LoadTodoItems() ([]models.TodoItem, error) {
	data, err := os.ReadFile(fileName)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.TodoItem{}, nil
		}
		return nil, err
	}
	var todos []models.TodoItem
	if err := json.Unmarshal(data, &todos); err != nil {
		return nil, err
	}
	return todos, nil
}

// store todo items in json file
func SaveTodoItems(todos []models.TodoItem) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(fileName, data, 0644)
}

// returns the index of the next todo item in the list
func GetNextId(items []models.TodoItem) int {
	maxId := 0
	for _, item := range items {
		if item.Id > maxId {
			maxId = item.Id
		}
	}
	return maxId + 1
}
