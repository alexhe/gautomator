package flue

// This is a basic example
// Thanks http://thenewstack.io/make-a-restful-json-api-go/ for the tutorial
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"time"
)

type TaskTest struct {
	Name      string
	Completed bool
	Due       time.Time
}

type Tasks []TaskTest

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}

func TaskIndex(w http.ResponseWriter, r *http.Request) {
	tasks := Tasks{
		TaskTest{Name: "Write presentation"},
		TaskTest{Name: "Host meetup"},
	}

	json.NewEncoder(w).Encode(tasks)
}

func TaskShow(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["Name"]
	fmt.Fprintln(w, "Task show:", name)
}
