package flue

import (
	"github.com/gonum/matrix/mat64"
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
func Runner(taskStructure *TaskGraphStructure, taskId int, taskStructureChan <-chan *TaskGraphStructure, doneChan chan<- int, wg *sync.WaitGroup) {
	return
	task := taskStructure.Tasks[taskId]
	log.Printf("[Id:%v, Name:%v] Queued", taskId, task.Name)
	for {
		loop := 1
		log.Printf("[Id:%v, Name:%v] loop %v", taskId, task.Name, loop)
		// Let's go unless we cannot
		loop += 1
		letsGo := true
		// For each dependency of the task
		// We can run if the sum of the element of the column Id of the current task is 0
		var taskDeps *mat64.Vector
		taskDeps = taskStructure.AdjacencyMatrix.ColView(taskId)
		for v := range taskDeps.RawVector().Data {
			log.Printf("[Id:%v] src: %v", taskId, v)
			if v != 0 {
				letsGo = false
			}
		}
		if letsGo == true {
			proto := "tcp"
			socket := "localhost:4546"
			sleepTime := random(5, 15)
			task.Module = "sleep"
			task.Args = []string{strconv.Itoa(sleepTime)}
			task.Status = -1
			log.Printf("[%v] Running (%v %v)", task.Name, task.Module, task.Args[0])
			log.Printf("[%v] Connecting in %v on %v", task.Name, proto, socket)
			task.StartTime = time.Now()
			task.Status = Client(task, &proto, &socket)
			task.EndTime = time.Now()
			// ... Do a lot of stufs...
			//time.Sleep(time.Duration(sleepTime) * time.Second)
			// Adjust the Status
			//task.Status = 2
			// Send it on the channel
			log.Printf("[%v] Done", task.Name)
			doneChan <- taskId
			wg.Done()
			return
		}
	}
}

// The advertize goroutine, reads the tasks from doneChannel and write the TaskGraphStructure back to the taskStructureChan
func Advertize(taskStructureChan chan<- *TaskGraphStructure, doneChan <-chan int) {
	for {
		<-doneChan
		//log.Printf("[%v] Finished, advertizing", doneTask.Name)
		//taskStructureChan <- taskGraphStructure
	}
}
