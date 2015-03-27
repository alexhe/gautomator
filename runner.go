package flue

import(
    "time"
    "log"
)

// Ths runner goroutine is a goroutinei which:
// Consume the TaskGraphStructure from a channel
// run the task given as arguments if the deps are done
// Post the task to the doncChannel once done
// 
func Runner(task *Task, taskStructureChan <-chan *TaskGraphStructure, doneChan chan<- *Task) {
    for { 
	taskStructure := <-taskStructureChan
	// Let's go unless we cannot
	letsGo := true
	// For each dependency of the task
	for _, dep := range task.Deps {
	    depTask := GetTask(dep, taskStructure) 
	    if depTask.Status != 2 {
		letsGo =false
	    }
	}
	if letsGo == true {
	    log.Printf("[%v] Running",task.Name)
	    // ... Do a lot of stufs...
	    time.Sleep(2 * time.Second)
	    // Adjust the Status
	    task.Status = 2
	    // Send it on the channel
	    doneChan <- task
	} else {
	    log.Printf("[%v] Waiting for deps")
	    for _, dep := range task.Deps {
		log.Printf("[%v] => %v",task.Name, dep)
	    }
	}
    }
}

// The advertize goroutine, reads the tasks from doneChannel and write the TaskGraphStructure back to the taskStructureChan
func Advertize(task *Task, initialTaskGraphStructure *TaskGraphStructure, tastStructureChan chan<- *TaskGraphStructure, doneChan <-chan *Task) {
    for {
	doneTask := <- doneChan
	log.Printf("DoneTask", doneTask.Name)
    }
}
