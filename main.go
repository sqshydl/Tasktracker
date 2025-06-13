package main

import (
	"time"
)

type Task struct {
	AutoId      int
	Description string
	Status      int
	CreatedAt   time.Time
	UpdateAt    time.Time
}

var currentId int = 0

func getNextId() int {
	currentId++
	return currentId
}

func add(newTask string) {

}

func main() {

}
