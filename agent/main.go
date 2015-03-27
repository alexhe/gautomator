package main

import (
//	"fmt"
//	"github.com/nu7hatch/gouuid"
	"github.com/owulveryck/flue"
	"log"
)

func main() {
	// Testing the DOT parsing...
	//	topologyDot, err := ioutil.ReadFile("test.dot")
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
	// Will have;
	// start waits for nothing
	// purge waits for start
	// installProduit1 waits for purge
	// installProduit2 waits for purge
	// startAll waits for installProduit1 AND instannProduit2
	// end waits for startAll
	log.Println("Parsing...")
	taskStructure := flue.ParseTasks(topologyDot)
	for i, task := range taskStructure.Tasks {
	    if task != nil {
		log.Printf("=> Tache %v: %s",i, task.Name)
		for j, dep := range task.Deps {
		    log.Printf("==> Deps[%v]: %v",j,dep)
		}
	    }
	}
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
	   	   wg.Wait()
	   	   fmt.Println("Done")
	*/
}
