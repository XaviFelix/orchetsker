package main

import (
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

// using a basic array for now
type TaskStore struct {
	Tasks []Task
}

const taskFile = "tasks.json"

// Pass a filename
func LoadTasks() (TaskStore, error) {
	store := TaskStore{}

	// check if file exists
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		return store, nil
	}

	// open file
	file, err := os.Open(taskFile)
	if err != nil {
		return store, fmt.Errorf("opening file: %w", err)
	}
	defer file.Close()

	// Check if file is empty
	stat, err := file.Stat()
	if err != nil {
		return store, fmt.Errorf("checking file stat: $%w", err)
	}
	if stat.Size() == 0 {
		return store, nil
	}

	// Everything works out, but checking if something went wrong with decoding
	err = json.NewDecoder(file).Decode(&store)
	if err != nil {
		return store, fmt.Errorf("decoding JSON: %w", err)
	}
	return store, nil
}

func SaveTasks(store TaskStore) error {
	file, err := os.Create(taskFile)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")

	return encoder.Encode(store)
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
	c.DeleteID = c.DeleteCmd.Int("id", 0, "ID of taks to delete")

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
		ListTasks(commands, store)
	// TODO: Done
	// TODO: Delete
	default:
		fmt.Println("Unknown command: ", os.Args[1])
		fmt.Println("Usage: <command> [task description]")
		fmt.Println("Available Commands: add, list, done, delete")
		os.Exit(1)
	}
}

// Adds a task to the json file
func AddTask(commands *Commands, store TaskStore) Task {

	// parse the description of the passed argument, located from args[2:]
	commands.Add.Parse(os.Args[2:])
	if len(commands.Add.Args()) == 0 {
		fmt.Println("Error: Provide a task description")
		os.Exit(1)
	}

	// gets the argument parsed from the cli
	description := commands.Add.Arg(0)

	// create a new id number based on array length
	newID := len(store.Tasks) + 1

	// append a new instance of a Task to the list of tasks
	store.Tasks = append(store.Tasks, Task{
		ID:          newID,
		Description: description,
		Done:        false,
	})
	fmt.Printf("Added task %d: %s\n", newID, description)
	// Save the task to json file
	SaveTasks(store)

	// return task to ensure it was created properly for debug purposes
	return store.Tasks[newID-1]
}

// list all tasks:
func ListTasks(commands *Commands, store TaskStore) []Task {
	commands.List.Parse(os.Args[2:])
	var stateOfTasks []Task

	if len(store.Tasks) == 0 {
		fmt.Println("No tasks yet")
		os.Exit(1) // changing into loop so instead, exit is temporary for now
	}

	for _, task := range store.Tasks {
		var status string
		if task.Done {
			status += "X"
		}
		stateOfTasks = append(stateOfTasks, task)
		fmt.Printf("[%s] %d: %s\n", status, task.ID, task.Description)
	}

	return stateOfTasks
}
