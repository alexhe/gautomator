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
	var nodesFileJson = flag.String("nodes", "", "json file for node definition")
	flag.Parse()
	dotFiles = flag.Args()
	log.Println("Server mode")
	proto := "tcp"
	socket := "localhost:4546"
	var wg sync.WaitGroup
	go gautomator.Rserver(&proto, &socket)
	wg.Add(1)

	if len(dotFiles) != 0 {
		log.Println("Client mode")

		// Parse the dot file, recompose and reconstruct a big graph
		taskStructure := gautomator.ParseDotFiles(dotFiles)
		// Parse the nodes.json and adapt the tasks
		taskInstances := gautomator.ParseNode(nodesFileJson)
		// Parse the instanciate the task with the information from the node file
		taskStructure.InstanciateTaskStructure(taskInstances)

		//taskStructure.PrintAdjacencyMatrix()
		// Entering the workers area
		doneChan := make(chan *gautomator.Task)

		// For each task, launch a goroutine
		for taskIndex, _ := range taskStructure.Tasks {
			go gautomator.Runner(taskStructure.Tasks[taskIndex], doneChan, &wg)
			wg.Add(1)
		}
		go gautomator.Advertize(taskStructure, doneChan)

		// This is the web displa
		router := gautomator.NewRouter(taskStructure)

		go log.Fatal(http.ListenAndServe(":8080", router))

		// Wait for all the runner(s) to be finished
	}
	wg.Wait()
}
