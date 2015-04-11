package gautomator

import (
	"github.com/gonum/matrix/mat64" // Matrix
	"log"
	"math/rand" // Temp
	//"strconv"
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
	log.Printf("[%v:%v] Queued", task.Id, task.Name)
	for {
		letsGo := <-task.TaskCanRunChan
		// For each dependency of the task
		// We can run if the sum of the element of the column Id of the current task is 0

		if letsGo == true {
			proto := "tcp"
			socket := task.Node
			//sleepTime := random(1, 5)
			// Stupid trick to make shell works... A Shell module will be implemented later"
			if task.Module == "shell" {
				task.Module = "echo"
				n := len(task.Args)
				task.Args[n] = "|"
				task.Args[n+1] = "/bin/ksh"
			}
			//task.Module = "sleep"
			//task.Args = []string{strconv.Itoa(sleepTime)}
			task.Status = -1
			log.Printf("[%v:%v] Running (%v %v)", task.Id, task.Name, task.Module, task.Args[0])
			//log.Printf("[%v] Connecting in %v on %v", task.Name, proto, socket)
			task.StartTime = time.Now()
			if task.Module != "dummy" && task.Module != "meta" {
				task.Status = Client(task, &proto, &socket)
			} else {
				task.Status = 0
			}
			task.EndTime = time.Now()
			// ... Do a lot of stufs...
			//time.Sleep(time.Duration(sleepTime) * time.Second)
			// Adjust the Status
			//task.Status = 2
			// Send it on the channel
			log.Printf("[%v:%v] Done", task.Id, task.Name)
			doneChan <- task
			wg.Done()
			return
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
	doneAdjacency := mat64.DenseCopyOf(taskStructure.AdjacencyMatrix)
	for {
		task := <-doneChan

		//TODO : There is absolutely no need to change the adjacency matrix AT ALL
		// It is a lot better to create a copy

		// Adapting the Adjacency matrix...
		// TaskId is finished, it cannot be the source of any task anymore
		// Set the row at 0
		rowSize, colSize := doneAdjacency.Dims()
		for c := 0; c < colSize; c++ {
			doneAdjacency.Set(task.Id, c, float64(0))
		}
		// For each dependency of the task
		// We can run if the sum of the element of the column Id of the current task is 0
		for taskIndex, _ := range taskStructure.Tasks {
			sum := float64(0)
			for r := 0; r < rowSize; r++ {
				sum += doneAdjacency.At(r, taskIndex)
			}

			if sum == 0 && taskStructure.Tasks[taskIndex].Status == -2 {
				taskStructure.Tasks[taskIndex].TaskCanRunChan <- true
			}
		}
	}
}
