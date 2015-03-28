package main

import (
	//	"fmt"
	//	"github.com/nu7hatch/gouuid"
	"github.com/owulveryck/flue"
	"log"
	"net/http"
	"sync"
)

func main() {
	// Testing the DOT parsing...
	//	topologyDot, err := ioutil.ReadFile("test.dot")
	/*
			topologyDot := []byte(`
		digraph layer3Tasks {
			start -> purge;
			test -> start;
			purge -> installProduit1;
			purge -> installProduit2;
			installProduit1 -> startAll;
			installProduit2 -> startAll;
			startAll -> end;

		}
		`)
	*/
	topologyDot := []byte(`
digraph layer3Tasks {
    A -> B;
    B -> D;
    B -> C;
    B -> E;
    D -> F;
    C -> F;
    F -> G;
    E -> G;
    G -> end;

    node [module = sleep,args = 2 3 4]; A;
}
`)
	log.Println("Parsing...")
	taskStructure := flue.ParseTasks(topologyDot)
	// How many tasks

	var wg sync.WaitGroup
	taskStructureChan := make(chan *flue.TaskGraphStructure)
	doneChan := make(chan *flue.Task)
	for i, task := range taskStructure.Tasks {
		if task != nil {
			go flue.Runner(taskStructure, task, taskStructureChan, doneChan, &wg)
			wg.Add(1)
			log.Printf("=> Tache %v: %s", i, task.Name)
			for j, dep := range task.Deps {
				log.Printf("==> Deps[%v]: %v", j, dep)
			}
		}
	}
	go flue.Advertize(taskStructure, taskStructureChan, doneChan)
	router := flue.NewRouter(*taskStructure)

	go log.Fatal(http.ListenAndServe(":8080", router))
	wg.Wait()
}
