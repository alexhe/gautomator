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
		for _, dotFile := range dotFiles {

			var topologyDot []byte
			topologyDot, err := ioutil.ReadFile(dotFile)
			if err != nil {
				log.Panic("Cannot read file")
			}

			log.Println("Parsing...")
			taskStructure := flue.ParseTasks(topologyDot)

			var wg sync.WaitGroup
			doneChan := make(chan *flue.Task)

			// For each task, if it can run, place true in its communication channel
			for taskIndex, _ := range taskStructure.Tasks {
				log.Printf("taskIndex: %v", taskIndex)
				go flue.Runner(taskStructure.Tasks[taskIndex], doneChan, &wg)
				wg.Add(1)
			}
			go flue.Advertize(taskStructure, doneChan)

			router := flue.NewRouter(taskStructure)

			go log.Fatal(http.ListenAndServe(":8080", router))
			wg.Wait()
		}
	}
}
