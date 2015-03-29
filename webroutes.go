package flue

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

// Courtesy of http://stackoverflow.com/questions/26211954/how-do-i-pass-arguments-to-my-handler
func showTasks(w http.ResponseWriter, r *http.Request, taskStruct TaskGraphStructure) {
	json.NewEncoder(w).Encode(taskStruct)
}

func NewRouter(taskStruct TaskGraphStructure) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")
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
