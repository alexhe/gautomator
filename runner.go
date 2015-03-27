package flue

import(
    "time"
    "log"
    "sync"
    "math/rand" // Temp
)



func random(min, max int) int {
    rand.Seed(time.Now().Unix())
    return rand.Intn(max - min) + min
}
// Ths runner goroutine is a goroutinei which:
// Consume the TaskGraphStructure from a channel
// run the task given as arguments if the deps are done
// Post the task to the doncChannel once done
// 
func Runner(taskStructure *TaskGraphStructure, task *Task, taskStructureChan <-chan *TaskGraphStructure, doneChan chan<- *Task, wg *sync.WaitGroup) {
    log.Printf("[%v] Queued",task.Name)
    for { 
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
	    time.Sleep(time.Duration(random(1, 16)) * time.Second)
	    // Adjust the Status
	    task.Status = 2
	    // Send it on the channel
	    log.Printf("[%v] Done", task.Name)
	    doneChan <- task
	    //	    log.Printf("[%v] channel: Done", task.Name)
	    wg.Done()
	    return
	}
    }
}

// The advertize goroutine, reads the tasks from doneChannel and write the TaskGraphStructure back to the taskStructureChan
func Advertize(taskGraphStructure *TaskGraphStructure, taskStructureChan chan<- *TaskGraphStructure, doneChan <-chan *Task) {
    for {
	<- doneChan
	//log.Printf("[%v] Finished, advertizing", doneTask.Name)
	//taskStructureChan <- taskGraphStructure
    }
}
