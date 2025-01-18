package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type TodoStruct struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
}

type TaskFile struct {
	Todos []TodoStruct `json:"todos"`
}

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}

func saveLocally(data TaskFile) {
	filePath := "./tasks.json"
	marshalJson, err := json.MarshalIndent(data, "", "  ")
	checkErr(err)
	err = os.WriteFile(filePath, marshalJson, 0644)
}

func main() {
	validStatuses := []string{"in-progress", "done", "todo"}
	filePath := "./tasks.json"

	taskFile, err := os.ReadFile(filePath)
	var data TaskFile

	if err != nil {
		data = TaskFile{Todos: []TodoStruct{}}

		jsonData, marshalErr := json.MarshalIndent(data, "", "  ")
		checkErr(marshalErr)

		err = os.WriteFile(filePath, jsonData, 0644)
		checkErr(err)
	} else {
		err = json.Unmarshal(taskFile, &data)
		checkErr(err)
	}

	idCounter := 0

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		splitText := strings.Split(line, " ")
		action := splitText[1]

		switch action {
		case "add":
			taskName := splitText[2]
			if len(data.Todos) > 0 {
				lastId := data.Todos[len(data.Todos)-1].Id
				newTodo := TodoStruct{
					Id:          lastId + 1,
					Description: taskName,
					Status:      "todo",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				data.Todos = append(data.Todos, newTodo)
			} else {
				newTodo := TodoStruct{
					Id:          idCounter + 1,
					Description: taskName,
					Status:      "todo",
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				data.Todos = append(data.Todos, newTodo)
			}
			saveLocally(data)
			fmt.Printf("New Task Added!\n -> task: %v / status: %v / ID: %v\n", taskName, "todo", data.Todos[len(data.Todos)-1].Id)
		case "update":
			idToUpdate := splitText[2]
			parsedId, err := strconv.Atoi(idToUpdate)
			checkErr(err)
			updatedDesc := splitText[3]
			var taskToUpdate TodoStruct
			for i := range data.Todos {
				if parsedId == data.Todos[i].Id {
					data.Todos[i].Description = updatedDesc
					taskToUpdate = data.Todos[i]
				}
			}
			saveLocally(data)
			fmt.Printf("Updated Task!\n -> task: %v / status: %v / ID: %v\n", updatedDesc, taskToUpdate.Status, parsedId)
		case "delete":
			idToDelete := splitText[2]
			parsedId, err := strconv.Atoi(idToDelete)
			checkErr(err)
			for i, v := range data.Todos {
				if v.Id == parsedId {
					// Delete the task by appending the parts before and after the element
					data.Todos = append(data.Todos[:i], data.Todos[i+1:]...)
					saveLocally(data)
					fmt.Printf("Task with ID %d deleted!\n", parsedId)
				}
			}
		case "list":
			if len(splitText) == 3 {
				byStatus := splitText[2]
				if contains(validStatuses, byStatus) {
					for i, v := range data.Todos {
						if v.Status == byStatus {
							fmt.Printf("%v) task: %v / status: %v / ID: %v\n", i+1, v.Description, v.Status, v.Id)
						}
					}
				} else {
					fmt.Printf("Enter a valid status!\n")
				}
			} else {
				for i, v := range data.Todos {
					fmt.Printf("%v) task: %v / status: %v / ID: %v\n", i+1, v.Description, v.Status, v.Id)
				}
			}
		case "mark-in-progress":
			idToUpdate := splitText[2]
			parsedId, err := strconv.Atoi(idToUpdate)
			checkErr(err)
			var taskToUpdate TodoStruct
			for i := range data.Todos {
				if parsedId == data.Todos[i].Id {
					data.Todos[i].Status = "in-progress"
					taskToUpdate = data.Todos[i]
				}
			}
			saveLocally(data)
			fmt.Printf("Status Updated!\n -> task: %v / status: %v / ID: %v\n", taskToUpdate.Description, taskToUpdate.Status, parsedId)
		case "mark-done":
			idToUpdate := splitText[2]
			parsedId, err := strconv.Atoi(idToUpdate)
			checkErr(err)
			var taskToUpdate TodoStruct
			for i := range data.Todos {
				if parsedId == data.Todos[i].Id {
					data.Todos[i].Status = "done"
					taskToUpdate = data.Todos[i]
				}
			}
			saveLocally(data)
			fmt.Printf("Status Updated!\n -> task: %v / status: %v / ID: %v\n", taskToUpdate.Description, taskToUpdate.Status, parsedId)
		case "exit":
			saveLocally(data)
			fmt.Printf("Your Todos are saved! See you soon.\n")
			return
		default:
			fmt.Println("Please enter a valid command!")
		}
	}

}
