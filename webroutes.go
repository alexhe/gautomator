package gautomator

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(taskStruct *TaskGraphStructure) *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Headers("Content-Type", "application/json", "X-Requested-With", "XMLHttpRequest")
	router.Methods("GET").Path("/svg").Name("SVG Representation").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
