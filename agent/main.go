package main

import (
	"flag"
	"github.com/owulveryck/gautomator"
	"log"
	"net/http"
	"sync"
)

func cleanup() {
	log.Println("cleanup")
}

func main() {

	//Parsing the dot
	var dotFiles []string
//	var nodesFileJson = flag.String("nodes", "", "json file for node definition")
	flag.Parse()
	dotFiles = flag.Args()
	if len(dotFiles) == 0 {
		log.Println("Server mode")
		proto := "tcp"
		socket := "localhost:4546"
		gautomator.Rserver(&proto, &socket)
	} else {
		log.Println("Client mode")

		taskStructure := gautomator.ParseDotFiles(dotFiles)
		// Parse the nodes.json and adapt the tasks
		//nodeStructure := gautomator.ParseNode(nodesFileJson)

		/*
		for _, node := range *nodeStructure {
		    log.Printf("taskName: %v, module: %v",node.Taskname,node.Module)
		}
		*/
		// Entering the workers area
		var wg sync.WaitGroup
		doneChan := make(chan *gautomator.Task)

		// For each task, if it can run, place true in its communication channel
		for taskIndex, _ := range taskStructure.Tasks {
			go gautomator.Runner(taskStructure.Tasks[taskIndex], doneChan, &wg)
			wg.Add(1)
		}
		go gautomator.Advertize(taskStructure, doneChan)

		// This is the web displa
		router := gautomator.NewRouter(taskStructure)

		go log.Fatal(http.ListenAndServe(":8080", router))

		// Wait for all the runner(s) to be finished
		wg.Wait()
	}
}
