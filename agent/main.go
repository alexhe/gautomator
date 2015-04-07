package main

import (
	"flag"
	"github.com/owulveryck/flue"
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
	//var nodesFileJson = flag.String("nodes", "", "json file for node definition")
	flag.Parse()
	dotFiles = flag.Args()
	if len(dotFiles) == 0 {
		log.Println("Server mode")
		proto := "tcp"
		socket := "localhost:4546"
		flue.Rserver(&proto, &socket)
	} else {
		log.Println("Client mode")

		taskStructure := flue.ParseDotFiles(dotFiles)
		// Parse the nodes.json and adapt the tasks

		// Entering the workers area
		var wg sync.WaitGroup
		doneChan := make(chan *flue.Task)

		// For each task, if it can run, place true in its communication channel
		for taskIndex, _ := range taskStructure.Tasks {
			go flue.Runner(taskStructure.Tasks[taskIndex], doneChan, &wg)
			wg.Add(1)
		}
		go flue.Advertize(taskStructure, doneChan)

		// This is the web displa
		router := flue.NewRouter(taskStructure)

		go log.Fatal(http.ListenAndServe(":8080", router))

		// Wait for all the runner(s) to be finished
		wg.Wait()
	}
}
