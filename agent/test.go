package main

import (
	"flag"
	"github.com/owulveryck/flue"
	"io/ioutil"
	"log"
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
	log.Println("Client mode")

	var topologyDot []byte
	topologyDot, err := ioutil.ReadFile(*dotFile)
	if err != nil {
		log.Panic("Cannot read file")
	}

	log.Println("Parsing...")
	flue.ParseTasks(topologyDot)
	// How many tasks
}
