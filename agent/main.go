package main

import (
	"encoding/json"
	"flag"
	"github.com/owulveryck/flue"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	//	"os"
	//	    "syscall"
	//	"os/signal"
)

func cleanup() {
	log.Println("cleanup")
}

func main() {

	/*
		// Catching the interrupt
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		signal.Notify(c, syscall.SIGTERM)
		go func() {
		    <-c
		    cleanup()
		    os.Exit(1)
		}()
	*/
	//Parsing the dot
	dotFile := flag.String("dot", "", "The dot file")

	flag.Parse()
	if *dotFile == "" {
		log.Println("Server mode")
		proto := "tcp"
		socket := "localhost:4546"
		flue.Rserver(&proto, &socket)
	} else {
		log.Println("Client mode")

		var topologyDot []byte
		topologyDot, err := ioutil.ReadFile(*dotFile)
		if err != nil {
			log.Panic("Cannot read file")
		}

		log.Println("Parsing...")
		taskStructure := flue.ParseTasks(topologyDot)

		// TEST
		sigmaStructure := flue.GetSigmaStructure(taskStructure)
		jsonOutput, _ := json.Marshal(sigmaStructure)
		log.Println(string(jsonOutput))
		// End TEST
		var wg sync.WaitGroup
		doneChan := make(chan *flue.Task)

		// END DEBUG

		// For each task, if it can run, place true in its communication channel
		for taskIndex, _ := range taskStructure.Tasks {
			log.Printf("taskIndex: %v", taskIndex)
			go flue.Runner(taskStructure.Tasks[taskIndex], doneChan, &wg)
			wg.Add(1)
		}
		go flue.Advertize(taskStructure, doneChan)

		router := flue.NewRouter(*taskStructure)

		go log.Fatal(http.ListenAndServe(":8080", router))
		wg.Wait()
	}
}
