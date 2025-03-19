package main

// TODO: create a flag for running in a loop or as one commmand

import (
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

// Maintains order
type TaskEntry struct {
	Task  Task
	Value *list.Element
}

// LinkedMap
type TaskStorage struct {
	Tasks map[int]*TaskEntry
	Order *list.List
}

// initalizes the linkedmap
func newTaskStorage() *TaskStorage {
	return &TaskStorage{
		Tasks: make(map[int]*TaskEntry),
		Order: list.New(),
	}
}

const taskFile = "tasks.json"

func LoadTasks() (*TaskStorage, error) {
	store := newTaskStorage()
	file, err := os.Open(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return store, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// load tasks from json into a slice
	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		log.Fatalf("Error decoding json: %v", err)
	}

	// use the slice to add entries to the linked map (store)
	for _, task := range tasks {
		entry := &TaskEntry{Task: task}
		entry.Value = store.Order.PushBack(entry)
		store.Tasks[task.ID] = entry
	}
	return store, nil
}

func SaveTask(store *TaskStorage) error {
	file, err := os.OpenFile("tasks.json", os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	tasks := make([]Task, 0, store.Order.Len())
	for e := store.Order.Front(); e != nil; e = e.Next() {
		tasks = append(tasks, e.Value.(*TaskEntry).Task)
	}

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	return encoder.Encode(tasks)
}

type Commands struct {
	Add       *flag.FlagSet
	List      *flag.FlagSet
	Done      *flag.FlagSet
	DoneID    *int
	DeleteCmd *flag.FlagSet
	DeleteID  *int
}

func NewCommands() *Commands {

	c := &Commands{
		Add:       flag.NewFlagSet("add", flag.ExitOnError),
		List:      flag.NewFlagSet("list", flag.ExitOnError),
		Done:      flag.NewFlagSet("done", flag.ExitOnError),
		DeleteCmd: flag.NewFlagSet("delete", flag.ExitOnError),
	}
	c.DoneID = c.Done.Int("id", 0, "ID of task to mark done")
	c.DeleteID = c.DeleteCmd.Int("id", 0, "ID of task to delete")

	return c
}

func main() {
	store, err := LoadTasks() // load task from json file
	commands := NewCommands() // instance of available commands

	if err != nil {
		fmt.Println("Error loading tasks: ", err)
		os.Exit(1)
	}

	// Arguments passed must be greater than 2
	if len(os.Args) < 2 {
		fmt.Println("Usage: <command> [task description]")
		fmt.Println("Commands: add, list, done, delete")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "add":
		AddTask(commands, store)
	case "list":
		ListTasks(commands, store)
	// TODO: Done
	// TODO: Delete
	default:
		fmt.Println("Unknown command: ", os.Args[1])
		fmt.Println("Usage: task <command> [args]")
		fmt.Println("Available Commands: add, list, done, delete")
		os.Exit(1)
	}
}

// Adds a task to the json file
func AddTask(commands *Commands, store *TaskStorage) Task {

	// parse the description of the passed argument, located from args[2:]
	commands.Add.Parse(os.Args[2:])
	if len(commands.Add.Args()) == 0 {
		fmt.Println("Error: Provide a task description")
		os.Exit(1)
	}

	// gets the argument parsed from the cli
	description := commands.Add.Arg(0)

	// create a new id number based on array length
	newID := len(store.Tasks) + 1 // Fix this

	// append a new instance of a Task to the list of tasks
	entry := &TaskEntry{Task: Task{ID: newID, Description: description, Done: false}}
	entry.Value = store.Order.PushBack(entry)
	store.Tasks[newID] = entry

	fmt.Printf("Added task %d: %s\n", newID, description)
	// Save the task to json file
	if err := SaveTask(store); err != nil {
		fmt.Println("Error saving tasks: ", err)
		os.Exit(1)
	}

	// return task to ensure it was created properly for debug purposes
	return store.Tasks[newID].Task
}

// list all tasks:
func ListTasks(commands *Commands, store *TaskStorage) []Task {
	commands.List.Parse(os.Args[2:])
	// var stateOfTasks []Task

	if len(store.Tasks) == 0 {
		fmt.Println("No tasks yet")
		os.Exit(1) // changing into loop so instead, exit is temporary for now
	}

	// Ordered
	tasks := make([]Task, 0, store.Order.Len())
	for e := store.Order.Front(); e != nil; e = e.Next() {
		var status string
		if e.Value.(*TaskEntry).Task.Done {
			status += "X"
		}
		tasks = append(tasks, e.Value.(*TaskEntry).Task)
		fmt.Printf("[%s] %d: %s\n",
			status, e.Value.(*TaskEntry).Task.ID,
			e.Value.(*TaskEntry).Task.Description)
	}

	// Unordered
	// for _, entry := range store.Tasks {
	// 	var status string
	// 	if entry.Task.Done {
	// 		status += "X"
	// 	}
	// 	stateOfTasks = append(stateOfTasks, entry.Task)
	// 	fmt.Printf("[%s] %d: %s\n", status, entry.Task.ID, entry.Task.Description)
	// }

	return tasks
}
