package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

type Task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      int       `json:"status"` //0 = Not yet, 1 = Progress, 2 = Done
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

var tasks []Task
var currentId int = 0

const filename = "task.json"

func getNextId() int {
	currentId++
	return currentId
}

// Load tasks from JSON file
func loadTasks() error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// File doesn't exist, start with empty tasks
		tasks = []Task{}
		return nil
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	if len(data) == 0 {
		// File is empty, start with empty tasks
		tasks = []Task{}
		return nil
	}

	err = json.Unmarshal(data, &tasks)
	if err != nil {
		return fmt.Errorf("error unmarshaling JSON: %v", err)
	}

	// Update currentId to be the highest existing ID
	for _, task := range tasks {
		if task.Id > currentId {
			currentId = task.Id
		}
	}

	return nil
}

// Save tasks to JSON file
func saveTasks() error {
	data, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %v", err)
	}

	err = ioutil.WriteFile(filename, data, 0644)
	if err != nil {
		return fmt.Errorf("error writing file: %v", err)
	}

	fmt.Printf("Tasks saved to %s\n", filename)
	return nil
}

func add() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Insert New Task: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)

	newTask := Task{
		Id:          getNextId(),
		Description: desc,
		Status:      0,
		CreatedAt:   time.Now(),
		UpdateAt:    time.Now(),
	}

	// Add the new task to the tasks slice
	tasks = append(tasks, newTask)

	// Save to file
	err := saveTasks()
	if err != nil {
		fmt.Printf("Error saving task: %v\n", err)
		return
	}

	fmt.Printf("Task added successfully! ID: %d\n", newTask.Id)
}

// Function to display all tasks
func listTasks() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	fmt.Println("\n--- All Tasks ---")
	for _, task := range tasks {
		status := "Not Started"
		if task.Status == 1 {
			status = "In Progress"
		} else if task.Status == 2 {
			status = "Done"
		}

		fmt.Printf("ID: %d | Status: %s | Description: %s\n",
			task.Id, status, task.Description)
		fmt.Printf("Created: %s | Updated: %s\n\n",
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			task.UpdateAt.Format("2006-01-02 15:04:05"))
	}
}

func main() {
	// Load existing tasks from file
	err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	fmt.Println("Task Tracker")
	fmt.Println("Commands: add, list, quit")

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\nEnter command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(strings.ToLower(command))

		switch command {
		case "add":
			add()
		case "list":
			listTasks()
		case "quit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Use: add, list, or quit")
		}
	}
}
