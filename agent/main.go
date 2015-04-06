package main

import (
	"flag"
	"github.com/owulveryck/flue"
	"io/ioutil"
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
	flag.Parse()
	dotFiles = flag.Args()
	if len(dotFiles) == 0 {
		log.Println("Server mode")
		proto := "tcp"
		socket := "localhost:4546"
		flue.Rserver(&proto, &socket)
	} else {
		log.Println("Client mode")
		// Parsing each dot file
		var taskStructureArray []*flue.TaskGraphStructure

		taskStructureArray = make([]*flue.TaskGraphStructure, len(dotFiles), len(dotFiles))
		for index, dotFile := range dotFiles {

			var topologyDot []byte
			topologyDot, err := ioutil.ReadFile(dotFile)
			if err != nil {
				log.Panic("Cannot read file: ", dotFile)
			}

			log.Printf("Parsing the file %v (%v)...", dotFile, index)
			taskStructureArray[index] = flue.ParseTasks(topologyDot)
			log.Printf("taskStructureArray[%v] filled", index)
		}
		var taskStructure *flue.TaskGraphStructure
		taskStructure = nil
		for index, taskStruct := range taskStructureArray {
			if index == 0 {
				taskStructure = taskStruct
			} else {
				taskStructure = taskStructure.AugmentTaskStructure(taskStruct)
			}
		}
		//taskStructure.PrintAdjacencyMatrix()
		// Entering the workers area
		var wg sync.WaitGroup
		doneChan := make(chan *flue.Task)

		// For each task, if it can run, place true in its communication channel
		for taskIndex, _ := range taskStructure.Tasks {
			log.Printf("taskIndex: %v", taskIndex)
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
