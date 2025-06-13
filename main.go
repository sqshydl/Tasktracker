package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

type Task struct {
	AutoId      int       `json:"id"`
	Description string    `json:"description"`
	Status      int       `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdateAt    time.Time `json:"updated_at"`
}

var currentId int = 0

func getNextId() int {
	currentId++
	return currentId
}

func add() Task {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Insert New Task: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)

}

func main() {

}
