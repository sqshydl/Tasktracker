package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
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

// Helper function to find task by ID
func findTaskById(id int) (*Task, int) {
	for i, task := range tasks {
		if task.Id == id {
			return &tasks[i], i
		}
	}
	return nil, -1
}

// Helper function to get status string
func getStatusString(status int) string {
	switch status {
	case 0:
		return "Not Started"
	case 1:
		return "In Progress"
	case 2:
		return "Done"
	default:
		return "Unknown"
	}
}

// Helper function to get user input as integer
func getIntInput(prompt string) (int, error) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	return strconv.Atoi(input)
}

// Helper function to get user input as string
func getStringInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func add() {
	desc := getStringInput("Insert New Task: ")
	if desc == "" {
		fmt.Println("Task description cannot be empty.")
		return
	}

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
		fmt.Printf("ID: %d | Status: %s | Description: %s\n",
			task.Id, getStatusString(task.Status), task.Description)
		fmt.Printf("Created: %s | Updated: %s\n\n",
			task.CreatedAt.Format("2006-01-02 15:04:05"),
			task.UpdateAt.Format("2006-01-02 15:04:05"))
	}
}

// Function to view a specific task by ID
func viewTask() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	id, err := getIntInput("Enter task ID to view: ")
	if err != nil {
		fmt.Println("Invalid ID format. Please enter a number.")
		return
	}

	task, _ := findTaskById(id)
	if task == nil {
		fmt.Printf("Task with ID %d not found.\n", id)
		return
	}

	fmt.Println("\n--- Task Details ---")
	fmt.Printf("ID: %d\n", task.Id)
	fmt.Printf("Description: %s\n", task.Description)
	fmt.Printf("Status: %s\n", getStatusString(task.Status))
	fmt.Printf("Created: %s\n", task.CreatedAt.Format("2006-01-02 15:04:05"))
	fmt.Printf("Updated: %s\n", task.UpdateAt.Format("2006-01-02 15:04:05"))
}

// Function to edit a task
func editTask() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	id, err := getIntInput("Enter task ID to edit: ")
	if err != nil {
		fmt.Println("Invalid ID format. Please enter a number.")
		return
	}

	task, _ := findTaskById(id)
	if task == nil {
		fmt.Printf("Task with ID %d not found.\n", id)
		return
	}

	fmt.Printf("\nCurrent task: %s\n", task.Description)
	fmt.Printf("Current status: %s\n", getStatusString(task.Status))

	fmt.Println("\nWhat would you like to edit?")
	fmt.Println("1. Description")
	fmt.Println("2. Status")
	fmt.Println("3. Both")

	choice, err := getIntInput("Enter choice (1-3): ")
	if err != nil {
		fmt.Println("Invalid choice. Please enter a number.")
		return
	}

	switch choice {
	case 1:
		newDesc := getStringInput("Enter new description: ")
		if newDesc != "" {
			task.Description = newDesc
			task.UpdateAt = time.Now()
		}
	case 2:
		fmt.Println("\nStatus options:")
		fmt.Println("0. Not Started")
		fmt.Println("1. In Progress")
		fmt.Println("2. Done")

		newStatus, err := getIntInput("Enter new status (0-2): ")
		if err != nil || newStatus < 0 || newStatus > 2 {
			fmt.Println("Invalid status. Please enter 0, 1, or 2.")
			return
		}
		task.Status = newStatus
		task.UpdateAt = time.Now()
	case 3:
		newDesc := getStringInput("Enter new description: ")
		if newDesc != "" {
			task.Description = newDesc
		}

		fmt.Println("\nStatus options:")
		fmt.Println("0. Not Started")
		fmt.Println("1. In Progress")
		fmt.Println("2. Done")

		newStatus, err := getIntInput("Enter new status (0-2): ")
		if err != nil || newStatus < 0 || newStatus > 2 {
			fmt.Println("Invalid status. Please enter 0, 1, or 2.")
			return
		}
		task.Status = newStatus
		task.UpdateAt = time.Now()
	default:
		fmt.Println("Invalid choice.")
		return
	}

	// Save to file
	err = saveTasks()
	if err != nil {
		fmt.Printf("Error saving task: %v\n", err)
		return
	}

	fmt.Printf("Task %d updated successfully!\n", id)
}

// Function to delete a task
func deleteTask() {
	if len(tasks) == 0 {
		fmt.Println("No tasks found.")
		return
	}

	id, err := getIntInput("Enter task ID to delete: ")
	if err != nil {
		fmt.Println("Invalid ID format. Please enter a number.")
		return
	}

	task, index := findTaskById(id)
	if task == nil {
		fmt.Printf("Task with ID %d not found.\n", id)
		return
	}

	// Show task details before deletion
	fmt.Printf("\nTask to delete: %s\n", task.Description)
	fmt.Printf("Status: %s\n", getStatusString(task.Status))

	confirm := getStringInput("Are you sure you want to delete this task? (y/N): ")
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "yes" {
		fmt.Println("Deletion cancelled.")
		return
	}

	// Remove task from slice
	tasks = append(tasks[:index], tasks[index+1:]...)

	// Save to file
	err = saveTasks()
	if err != nil {
		fmt.Printf("Error saving tasks: %v\n", err)
		return
	}

	fmt.Printf("Task %d deleted successfully!\n", id)
}

func main() {
	// Load existing tasks from file
	err := loadTasks()
	if err != nil {
		fmt.Printf("Error loading tasks: %v\n", err)
		return
	}

	fmt.Println("Task Tracker")
	fmt.Println("Commands: add, list, view, edit, delete, quit")

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
		case "view":
			viewTask()
		case "edit":
			editTask()
		case "delete":
			deleteTask()
		case "quit":
			fmt.Println("Goodbye!")
			return
		default:
			fmt.Println("Unknown command. Use: add, list, view, edit, delete, or quit")
		}
	}
}
