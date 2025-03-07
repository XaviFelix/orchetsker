package main

import (
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	Done        bool   `json:"done"`
}

type TaskEntry struct {
	Task  Task
	Value *list.Element
}

// Changed to Map
type TaskStorage struct {
	Tasks map[int]*TaskEntry
	Order *list.List
}

func NewTaskStorage() *TaskStorage {
	return &TaskStorage{
		Tasks: make(map[int]*TaskEntry),
		Order: list.New(),
	}
}

const taskFile = "tasks.json"

// Pass a filename
func LoadTasks() (*TaskStorage, error) {
	store := NewTaskStorage()
	file, err := os.Open(taskFile)
	if err != nil {
		if os.IsNotExist(err) {
			return store, nil
		}
		return store, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	raw, err := os.ReadFile(taskFile)
	if err != nil {
		return store, fmt.Errorf("reading file: %w", err)
	}
	if len(raw) == 0 {
		return store, nil
	}
	fmt.Printf("Debug: tasks.json content: %s\n", raw)

	var tasks []Task
	if err := json.Unmarshal(raw, &tasks); err != nil {
		return store, fmt.Errorf("decoding JSON: %w (content: %s)", err, raw)
	}

	// Populate map and list
	for _, task := range tasks {
		entry := &TaskEntry{Task: task}
		entry.Value = store.Order.PushBack(entry)
		store.Tasks[task.ID] = entry
	}
	return store, nil
}

func SaveTasks(store *TaskStorage) error {
	file, err := os.Create(taskFile)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
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
		currentTask := AddTask(commands, store)
		fmt.Printf("Adding a task: %v", currentTask)
	case "list":
		// ListTasks(commands, store)
		fmt.Println("Still need to list tasks here")
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
	if err := SaveTasks(store); err != nil {
		fmt.Println("Error saving tasks: ", err)
		os.Exit(1)
	}

	// return task to ensure it was created properly for debug purposes
	return store.Tasks[newID].Task
}

// list all tasks:
// func ListTasks(commands *Commands, store *TaskStorage) []Task {
// 	commands.List.Parse(os.Args[2:])
// 	var stateOfTasks []Task

// 	if len(store.Tasks) == 0 {
// 		fmt.Println("No tasks yet")
// 		os.Exit(1) // changing into loop so instead, exit is temporary for now
// 	}

// 	for _, task := range store.Tasks {
// 		var status string
// 		if task.Done {
// 			status += "X"
// 		}
// 		stateOfTasks = append(stateOfTasks, task)
// 		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Description)
// 	}

// 	return stateOfTasks
// }
