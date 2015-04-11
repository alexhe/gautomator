package gautomator

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

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
	sigmaStructure := GetSigmaStructure(taskStructure)
	//	jsonOutput, _ := json.Marshal(sigmaStructure)
	json.NewEncoder(w).Encode(sigmaStructure)
}

func NewRouter(taskStruct *TaskGraphStructure) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")
	router.Methods("GET").Path("/svg").Name("TaskIndex").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		displaySvg(w, r, taskStruct)
	})
	router.Methods("GET").Path("/tasks").Name("TaskIndex").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		showTasks(w, r, taskStruct)
	})
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("../htdocs/static/")))

	/*
		for _, route := range routes {
			router.
				Methods(route.Method).
				Path(route.Pattern).
				Name(route.Name).
				Handler(route.HandlerFunc)
		}
	*/
	return router
}

/*
var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
		Route{
			"TaskIndex",
			"GET",
			"/tasks",
			showTasks,
		},
	Route{
		"TaskShow",
		"GET",
		"/tasks/{Name}",
		TaskShow,
	},
}
*/
