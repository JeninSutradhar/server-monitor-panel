package api

import (
	"errors"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Task struct { // public data type structure (Pascal cases on initial names)

	ID          int       `json:"id"`
	Description string    `json:"description"`
	CreatedTime time.Time `json:"created_time"`
	RunTime     time.Time `json:"run_time"` // for specific scheduling operation
	IsFinished  bool      `json:"is_finished"`
}

// general type to access in all place (Pascal cases on initial names)
type TaskStore struct { // also type to make type exports and accessible. struct name from `tasks` are declared on public
	tasks map[int]*Task

	sync.RWMutex
}

var tasks *TaskStore // and a struct implementation of exported variables .

// initializes an empty list of task struct to keep running information (its a in memory volatile solution) ( export )
func InitTasks() {
	tasks = &TaskStore{
		tasks: make(map[int]*Task),
	}
}

// gets the task available on that current struct ( export a method struct )
func ListTasks() map[int]*Task {

	tasks.RLock()
	defer tasks.RUnlock()

	copyTasks := make(map[int]*Task)

	for k, v := range tasks.tasks {

		copyTasks[k] = v
	}
	return copyTasks

}

// function for task submission that schedules when it should start to perform specific actions for each taks using a inmemory structure of type task( export  implement also method by Pascal case!)
func SubmitTask(desc string, runTime time.Time) (Task, error) {

	tasks.Lock()

	defer tasks.Unlock()
	newID := len(tasks.tasks)

	task := &Task{
		ID:          newID,
		Description: desc,
		CreatedTime: time.Now(),

		RunTime:    runTime,
		IsFinished: false,
	}

	tasks.tasks[newID] = task // all types declared on public level if the used methods , struct data implementation, this avoid to those ""type or variables by compiler". It can now see!.
	go executeTask(task)      // for  routine also  (pointer struct) implementation if exist in time method

	return *task, nil

}

// execute simulated task and do changes and return to list that its currently executing task with given parameter for simulating its current state with running for a duration(lower cases since local used no public by types struct at logic ).  Since they does  use that struct but not expose out at `tasks.go`,.
func executeTask(task *Task) {

	rand.Seed(rand.Int63())
	sec := rand.Intn(10) + 1

	<-time.After(time.Duration(sec) * time.Second)

	tasks.Lock()
	defer tasks.Unlock()
	task.IsFinished = true

	fmt.Printf("The  taks %d   description was:   %s, with new status now!.\n", task.ID, task.Description)

}

// gets specific task information with the specified ID(export types struct with names (must by pascal cases on methods).  )
func GetTask(taskID int) (Task, error) {

	tasks.RLock()

	defer tasks.RUnlock()

	task, ok := tasks.tasks[taskID]
	if ok {

		return *task, nil
	}
	return Task{}, errors.New(fmt.Sprintf("No tasks was found from list  by id %d. from routes", taskID))

}

// remove existing task(also a public implemented method with public  by a pascal cases on types). also that implementation do a error handling to validate it!

func DeleteTask(taskId int) error {

	tasks.Lock()

	defer tasks.Unlock()

	_, ok := tasks.tasks[taskId]
	if !ok {

		return fmt.Errorf("the tasks :%v is not implemented with this methods.", taskId)
	}

	delete(tasks.tasks, taskId)

	return nil

}
