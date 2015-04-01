package flue

import (
	"log"
	"math/rand" // Temp
	"strconv"
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
func Runner(task *Task, doneChan chan<- *Task, wg *sync.WaitGroup) {
	log.Printf("[Name:%v] Queued", task.Name)
	for {
		letsGo := <-task.TaskCanRunChan
		// For each dependency of the task
		// We can run if the sum of the element of the column Id of the current task is 0

		if letsGo == true {
			proto := "tcp"
			socket := "localhost:4546"
			sleepTime := random(1, 5)
			task.Module = "sleep"
			task.Args = []string{strconv.Itoa(sleepTime)}
			task.Status = -1
			log.Printf("[%v] Running (%v %v)", task.Name, task.Module, task.Args[0])
			//log.Printf("[%v] Connecting in %v on %v", task.Name, proto, socket)
			task.StartTime = time.Now()
			task.Status = Client(task, &proto, &socket)
			task.EndTime = time.Now()
			// ... Do a lot of stufs...
			//time.Sleep(time.Duration(sleepTime) * time.Second)
			// Adjust the Status
			//task.Status = 2
			// Send it on the channel
			log.Printf("[%v] Done", task.Name)
			doneChan <- task
			wg.Done()
			//return
		}
	}
}

// The advertize goroutine, reads the tasks from doneChannel and write the TaskGraphStructure back to the taskStructureChan
func Advertize(taskStructure *TaskGraphStructure, doneChan <-chan *Task) {
	// Let's launch the task that can initially run
	rowSize, _ := taskStructure.AdjacencyMatrix.Dims()
	for taskIndex, _ := range taskStructure.Tasks {
		sum := float64(0)
		for r := 0; r < rowSize; r++ {
			sum += taskStructure.AdjacencyMatrix.At(r, taskIndex)
		}
		if sum == 0 && taskStructure.Tasks[taskIndex].Status < 0 {
			taskStructure.Tasks[taskIndex].TaskCanRunChan <- true
		}
	}
	for {
		task := <-doneChan

		// Adapting the Adjacency matrix...
		// TaskId is finished, it cannot be the source of any task anymore
		// Set the row at 0
		rowSize, colSize := taskStructure.AdjacencyMatrix.Dims()
		for c := 0; c < colSize; c++ {
			taskStructure.AdjacencyMatrix.Set(task.Id, c, float64(0))
		}
		// For each dependency of the task
		// We can run if the sum of the element of the column Id of the current task is 0
		for taskIndex, _ := range taskStructure.Tasks {
			sum := float64(0)
			for r := 0; r < rowSize; r++ {
				sum += taskStructure.AdjacencyMatrix.At(r, taskIndex)
			}

			if sum == 0 && taskStructure.Tasks[taskIndex].Status == -2 {
				taskStructure.Tasks[taskIndex].TaskCanRunChan <- true
			}
		}
	}
}
