package gautomator

// This is a basic example
// Thanks http://thenewstack.io/make-a-restful-json-api-go/ for the tutorial
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func displaySvg(w http.ResponseWriter, r *http.Request, taskStructure *TaskGraphStructure) {
	subProcess := exec.Command("dot", "-Tsvg")

	stdin, err := subProcess.StdinPipe()
	if err != nil {
		fmt.Println(err) //replace with logger, or anything you want
	}
	defer stdin.Close() // the doc says subProcess.Wait will close it, but I'm not sure, so I kept this line

	subProcess.Stdout = w
	subProcess.Stderr = os.Stderr
	if err = subProcess.Start(); err != nil { //Use start, not run
		fmt.Println("An error occured: ", err) //replace with logger, or anything you want

	}
	taskStructure.PrintDot(stdin)
	// Command was successful
	stdin.Close()
	subProcess.Wait()

}

// Courtesy of http://stackoverflow.com/questions/26211954/how-do-i-pass-arguments-to-my-handler
func showTasks(w http.ResponseWriter, r *http.Request, taskStructure *TaskGraphStructure) {
	// A task is an action executed by a module
	type taskJson struct {
		Id     int      `json:"id"`
		Origin string   `json:"origin"`
		Name   string   `json:"name"` //the task name
		Node   string   `json:"node"` // The node name
		Module string   `json:"module"`
		Args   []string `json:"args"`
		Status int      `json:"status"` //-2: queued
		// -1: running
		// >=0 : return code
		StartTime time.Time `json:"startTime"`
		EndTime   time.Time `json:"endTime"`
	}
	var tasksJ []taskJson
	for _, task := range taskStructure.Tasks {
		var taskJ taskJson
		taskJ.Id = task.Id
		taskJ.Origin = task.Origin
		taskJ.Name = task.Name
		taskJ.Node = task.Node
		taskJ.Module = task.Module
		taskJ.Args = task.Args
		taskJ.Status = task.Status
		taskJ.StartTime = task.StartTime
		taskJ.EndTime = task.EndTime
		tasksJ = append(tasksJ, taskJ)
	}
	err := json.NewEncoder(w).Encode(tasksJ)
	if err != nil {
		fmt.Println(err)
	}
}
