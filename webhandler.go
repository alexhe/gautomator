package gautomator

// This is a basic example
// Thanks http://thenewstack.io/make-a-restful-json-api-go/ for the tutorial
import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
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
	//io.WriteString(stdin, "digraph G {\n")
	//io.WriteString(stdin, " a->b\n")
	//io.WriteString(stdin, "}\n")
	// Command was successful
	stdin.Close()
	subProcess.Wait()

}

// Courtesy of http://stackoverflow.com/questions/26211954/how-do-i-pass-arguments-to-my-handler
func showTasks(w http.ResponseWriter, r *http.Request, taskStructure *TaskGraphStructure) {
	//sigmaStructure := GetSigmaStructure(taskStructure)
	//	jsonOutput, _ := json.Marshal(sigmaStructure)
	json.NewEncoder(w).Encode(taskStructure.Tasks)
}
