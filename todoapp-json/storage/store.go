package storage

import (
	"encoding/json"
	"errors"
	"log/slog"
	"os"
	"todoapp-json/models"
)

type CommandType int

const (
	CmdAdd CommandType = iota
	CmdList
	CmdUpdate
	CmdDelete
)

type StoreCommand struct {
	Type       CommandType
	Item       models.TodoItem
	Id         int
	ResponseCh chan StoreResponse
}

type StoreResponse struct {
	TodoItems []models.TodoItem
	Err       error
}

var (
	StoreChan chan StoreCommand
	FilePath  string = "data/data.json" // Default, override from CLI
)

func StartStoreLoop() {
	todoItems, err := LoadTodoItems() //loading from jsonfile
	if err != nil {
		slog.Error("Failed to load todos", "err", err)
		todoItems = []models.TodoItem{}
	}

	for cmd := range StoreChan {
		switch cmd.Type {

		case CmdAdd:
			cmd.Item.Id = GetNextId(todoItems)
			todoItems = append(todoItems, cmd.Item)
			cmd.ResponseCh <- StoreResponse{Err: SaveTodoItems(todoItems)}

		case CmdList:
			cmd.ResponseCh <- StoreResponse{TodoItems: todoItems}

		case CmdUpdate:
			updated := false
			for i, item := range todoItems {
				if item.Id == cmd.Item.Id {
					todoItems[i] = cmd.Item
					updated = true
					break
				}
			}
			if !updated {
				cmd.ResponseCh <- StoreResponse{Err: errors.New("update : todo item not found")}
			} else {
				cmd.ResponseCh <- StoreResponse{Err: SaveTodoItems(todoItems)}
			}

		case CmdDelete:
			index := -1
			for i, item := range todoItems {
				if item.Id == cmd.Id {
					index = i
					break
				}
			}
			if index == -1 {
				cmd.ResponseCh <- StoreResponse{Err: errors.New("delete : todo item not found")}
				continue
			}
			todoItems = append(todoItems[:index], todoItems[index+1:]...)
			cmd.ResponseCh <- StoreResponse{Err: SaveTodoItems(todoItems)}
		}
	}
}

// get the next id after the last item
func GetNextId(todos []models.TodoItem) int {
	maxID := 0
	for _, item := range todos {
		if item.Id > maxID {
			maxID = item.Id
		}
	}
	return maxID + 1
}

// loading todo items from json file
func LoadTodoItems() ([]models.TodoItem, error) {
	file, err := os.Open(FilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []models.TodoItem{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var todos []models.TodoItem
	dec := json.NewDecoder(file)
	if err := dec.Decode(&todos); err != nil {
		return nil, errors.New("failed to decode JSON: " + err.Error())
	}

	return todos, nil
}

// store todo items to the json file
func SaveTodoItems(todos []models.TodoItem) error {
	data, err := json.MarshalIndent(todos, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(FilePath, data, 0644)
}

func SetDataFilePath(path string) {
	FilePath = path
}
