package main

// TODO: create a flag for running in a loop or as one commmand
// TODO: change the visibility of all of these methods
// TODO: Fix the duplicate task on 'done' bug, the problem lies in the SaveTask method
// TODO: Change the way i keep track of jsonl elements using an offset buffer

// TODO: finish the delete command

// TODO: Add a done date
//		 Then add a check in loadTasks.
//		 If a task that is done exceeds 10 days
// 		 delete it from the jsonl

import (
	"container/list"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

const taskFile = "tasks.jsonl"

// TODO: Converting to jsonl soon, refactor this code
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

	// var tasks []Task
	decoder := json.NewDecoder(file)
	for {
		var task Task
		if err := decoder.Decode(&task); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		entry := &TaskEntry{Task: task}
		entry.Value = store.Order.PushBack(entry)
		store.Tasks[task.ID] = entry
		// tasks = append(tasks, task)
	}
	return store, nil
}

// TODO: THis needs to save one task to json
// changing to jsonl soon, refactor this code
// add another argument here: task Task
func SaveTask(store *TaskStorage, task Task) error {
	file, err := os.OpenFile(taskFile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	return encoder.Encode(task)

	// TODO: modify these lines of code so that it encodes only one task
	// encoder := json.NewEncoder(file)
	// encoder.SetIndent("", " ")
	// return encoder.Encode(tasks)
}

type Commands struct {
	Add       *flag.FlagSet
	List      *flag.FlagSet
	Done      *flag.FlagSet
	DoneID    *int          // Use this to create a separate jsonl of done tasks
	DeleteCmd *flag.FlagSet // Use this to create a spearate jsonl of deleted tasks
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
	case "done":
		setDone(commands, store)
	// TODO: Delete
	case "delete":
		deleteTask(commands, store)
	default:
		fmt.Println("Unknown command: ", os.Args[1])
		fmt.Println("Usage: task <command> [args]")
		fmt.Println("Available Commands: add, list, done, delete")
		os.Exit(1)
	}
}

// This needs to delete from a json file
// NOTE: Figure out how to delete an entry in a jsonl file
// func deleteTask(commands *Commands, store *TaskStorage) {
// 	commands.DeleteCmd.Parse(os.Args[2:])
// 	if len(commands.DeleteCmd.Args()) == 0 {
// 		fmt.Println("Error: Provide a task id")
// 		os.Exit(1)
// 	}

// 	taskId, err := strconv.Atoi(commands.Done.Arg(0))
// 	if err != nil {
// 		log.Fatalf("Something went wrong converting taskID: %v", err)
// 	}

// 	// use the task id in order to find the element in the
// 	// store and in the jsonl

// }

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
	if err := SaveTask(store, entry.Task); err != nil {
		fmt.Println("Error saving tasks: ", err)
		os.Exit(1)
	}

	// return task to ensure it was created properly for debug purposes
	// TODO: return an error instead
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

// TODO: This saves a copy of a task that was set to done.
// NOTE: that means there is a duplicate: the original (not done) task
//
//	and the copy (is done)
func setDone(commands *Commands, store *TaskStorage) {
	if len(store.Tasks) == 0 {
		fmt.Println("No tasks yet")
		os.Exit(1) // changing into loop so instead, exit is temporary for now
	}

	commands.Done.Parse(os.Args[2:])

	taskId, err := strconv.Atoi(commands.Done.Arg(0))
	if err != nil {
		log.Fatalf("Something went wrong converting taskID: %v", err)
	}

	currentTask := store.Tasks[taskId]
	currentTask.Task.Done = true

	// TODO: Fix this
	// So what this does is that it creates a new
	// task on the jsonl file instead of updating it
	err = SaveTask(store, currentTask.Task) // pass the curretTask to this method
	if err != nil {
		log.Fatal("error saving current task: ", err)
	}
	ListTasks(commands, store)
}

// func deleteTask(commands *Commands, store *TaskStorage) {
// 	if len(store.Tasks) == 0 {
// 		fmt.Println("No taks yet")
// 		os.Exit(1)
// 	}

// 	commands.DeleteCmd.Parse(os.Args[2:])
// 	taskId, err := strconv.Atoi(commands.DeleteCmd.Arg(0))
// 	if err != nil {
// 		log.Fatalf("Something went wrong converting taskID: %v", err)
// 	}

// 	// This is deleted from the map but is not delted from the json file
// 	delete(store.Tasks, taskId)

// }
