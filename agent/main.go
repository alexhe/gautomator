package main

import (
//	"fmt"
//	"github.com/nu7hatch/gouuid"
	"github.com/owulveryck/flue"
	"sync"
	"log"
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

    node [shape = doublecircle];
}
`)
	log.Println("Parsing...")
	taskStructure := flue.ParseTasks(topologyDot)
	// How many tasks

	var wg sync.WaitGroup
	wg.Add(2)


	taskStructureChan := make(chan *flue.TaskGraphStructure)   	
	doneChan := make(chan *flue.Task)   	
	for i, task := range taskStructure.Tasks {
	    if task != nil {
		go flue.Runner(taskStructure, task, taskStructureChan, doneChan, &wg)
		log.Printf("=> Tache %v: %s",i, task.Name)
		for j, dep := range task.Deps {
		    log.Printf("==> Deps[%v]: %v",j,dep)
		}
	    }
	}
	go flue.Advertize(taskStructure,taskStructureChan, doneChan)

	/*
	   	if len(os.Args) < 2 {
	   		uuid, _ := uuid.NewV4()
	   		log.Println("We are a server, uuid is: ", string(uuid[:]))
	   		flue.Server("localhost:5678")
	   	} else {
	   		log.Println("We are a client...")
	   		myTasksChan := make(chan *flue.TopologyGraphStructure)
	   		for _, task := range myTasks.AllTheTasks {
	   			go flue.RunTask(task, myTasksChan, &wg)
	   		}
	   		myTasksChan <- myTasks
	   /*
	   		myTasksChan <- myTasks
	   		myTasksChan <- myTasks
	   		myTasksChan <- myTasks
	   		myTasksChan <- myTasks
	   	}
	   	   fmt.Println("Done")
	*/
       wg.Wait()

}
