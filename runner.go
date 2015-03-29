package flue

import (
	"log"
	"math/rand" // Temp
	"sync"
	"time"
)

func random(min int, max int) int {
	var bytes int
	bytes = min + rand.Intn(max)
	return int(bytes)
	//rand.Seed(time.Now().UTC().UnixNano())
	//return rand.Intn(max - min) + min
}

// Ths runner goroutine is a goroutinei which:
// Consume the TaskGraphStructure from a channel
// run the task given as arguments if the deps are done
// Post the task to the doncChannel once done
//
func Runner(taskStructure *TaskGraphStructure, task *Task, taskStructureChan <-chan *TaskGraphStructure, doneChan chan<- *Task, wg *sync.WaitGroup) {
	log.Printf("[%v] Queued", task.Name)
	for {
		// Let's go unless we cannot
		letsGo := true
		// For each dependency of the task
		for _, dep := range task.Deps {
			depTask := GetTask(dep, taskStructure)
			if depTask.Status != 2 {
				letsGo = false
			}
		}
		if letsGo == true {
			sleepTime := random(5, 15)
			log.Printf("[%v] Running (sleep for %v seconds)", task.Name, sleepTime)
			// ... Do a lot of stufs...
			time.Sleep(time.Duration(sleepTime) * time.Second)
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
		<-doneChan
		//log.Printf("[%v] Finished, advertizing", doneTask.Name)
		//taskStructureChan <- taskGraphStructure
	}
}
