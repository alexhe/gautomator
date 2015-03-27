package flue

import (
	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"sync"
	"syscall"
)

// RemoteCommand is the run parameters to be executed remotely
type RemoteCommandClient struct {
	Cmd        string
	Args       []string
	Stdin      io.Writer
	Stdout     io.Reader
	Stderr     io.Reader
	StatusChan libchan.Sender
}
type RemoteCommandServer struct {
	Cmd        string
	Args       []string
	Stdin      io.Reader
	Stdout     io.WriteCloser
	Stderr     io.WriteCloser
	StatusChan libchan.Sender
}

// CommandResponse is the returned response object from the remote execution
type CommandResponse struct {
	Status int
}

func Server(socketName string) {
	var listener net.Listener
	var err error
	//listener, err = net.Listen("unix", socketName)
	listener, err = net.Listen("tcp", socketName)
	defer os.Remove(socketName)
	if err != nil {
		log.Fatal(err)
	}

	tl, err := spdy.NewTransportListener(listener, spdy.NoAuthenticator)
	if err != nil {
		log.Fatal(err)
	}

	for {
		t, err := tl.AcceptTransport()
		if err != nil {
			log.Print(err)
			break
		}

		go func() {
			for {
				receiver, err := t.WaitReceiveChannel()
				if err != nil {
					log.Print(err)
					break
				}

				go func() {
					for {
						command := &RemoteCommandServer{}
						err := receiver.Receive(command)
						if err != nil {
							log.Print(err)
							break
						}

						log.Println("Will execute:", command.Cmd)
						cmd := exec.Command(command.Cmd, command.Args...)
						cmd.Stdout = command.Stdout
						cmd.Stderr = command.Stderr

						stdin, err := cmd.StdinPipe()
						if err != nil {
							log.Print(err)
							break
						}
						go func() {
							io.Copy(stdin, command.Stdin)
							stdin.Close()
						}()

						res := cmd.Run()
						command.Stdout.Close()
						command.Stderr.Close()
						returnResult := &CommandResponse{}
						if res != nil {
							if exiterr, ok := res.(*exec.ExitError); ok {
								returnResult.Status = exiterr.Sys().(syscall.WaitStatus).ExitStatus()
							} else {
								log.Print(res)
								returnResult.Status = 10
							}
						}

						log.Println("Sending back ", returnResult.Status)
						err = command.StatusChan.Send(returnResult)
						if err != nil {
							log.Print(err)
						}
					}
				}()
			}
		}()
	}
}

func RemoveTask(allTasksOriginal *TopologyGraphStructure, allTasks chan<- *TopologyGraphStructure, taskChan <-chan string, wg *sync.WaitGroup) {
	allTasks <- allTasksOriginal
	for {
		log.Println("RemoveTask consumming channel")
		taskLocal := <-taskChan
		// Remove the task from all the dependencies
		var temp []string
		temp = nil
		for _, atask := range allTasksOriginal.AllTheTasks {
			for _, deps := range allTasksOriginal.waiter[atask] {
				if deps != taskLocal {
					//log.Printf("Adding %s to waiter list %s", deps, atask)

					AppendTask(temp, deps)
				}
			}
			delete(allTasksOriginal.waiter, atask)
			if temp != nil {
				allTasksOriginal.waiter[atask] = temp
			}
			//DEBUG
			for _, entry := range allTasksOriginal.waiter[atask] {
				log.Println("DEBUG", entry)
			}
		}
		// Give it back to the channel
		//allTasks <- allTasksOriginal
		// log.Printf("[%s] Finished",task)
		log.Println("Writing back the channel")
		allTasks <- allTasksOriginal
		//allTasks <- allTasksOriginal
		//wg.Done()
	}
}

func RunTask(allTasksLocal *TopologyGraphStructure) {
	for _, task := range allTasksLocal.AllTheTasks {
		//allTasksLocal := <-allTasks
		//log.Println("Queuing:", task)
		//log.Println("DEBUG:", len(allTasksLocal.waiter[task]))
		if len(allTasksLocal.waiter[task]) != 0 {
			//if _, ok := allTasksLocal.waiter[task]; ok {
			log.Printf("[%s] Waiting for the folowing tasks to finish:", task)

			for _, s := range allTasksLocal.waiter[task] {
				log.Printf("%s ", s)
			}
			log.Printf("\n")
		} else {
			log.Printf("[%s] Let's go\n", task)
			command := &RemoteCommandClient{
				Cmd:    "date",
				Args:   []string{""},
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
				//StatusChan: remoteSender,
			}
			//Client(command, "/tmp/mysocket.sock")
			Client(command, "localhost:5678")
			log.Printf("[%s] Finished", task)

			var temp []string
			for _, atask := range allTasksLocal.AllTheTasks {
				for _, deps := range allTasksLocal.waiter[atask] {
					if deps != task {
						//log.Printf("Adding %s to waiter list %s", deps, atask)
						AppendTask(temp, deps)
					}
				}
				delete(allTasksLocal.waiter, atask)
				allTasksLocal.waiter[atask] = temp
				//DEBUG
				for _, entry := range allTasksLocal.waiter[atask] {
					log.Println("DEBUG", entry)
				}
			}
		}
	}
}

/*
func RunTask(task string, allTasks <-chan *TopologyGraphStructure, taskChan chan<- string) *TopologyGraphStructure {
	for {
		allTasksLocal := <-allTasks
		//log.Println("Queuing:", task)
		if _, ok := allTasksLocal.waiter[task]; ok {
			log.Printf("[%s] Waiting for the folowing tasks to finish:", task)

			for _, s := range allTasksLocal.waiter[task] {
				log.Printf("%s ", s)
			}
			log.Printf("\n")
		} else {
			log.Printf("[%s] Let's go", task)
			command := &RemoteCommandClient{
				Cmd:    "touch",
				Args:   []string{task},
				Stdin:  os.Stdin,
				Stdout: os.Stdout,
				Stderr: os.Stderr,
				//StatusChan: remoteSender,
			}
			//Client(command, "/tmp/mysocket.sock")
			Client(command, "localhost:5678")
			log.Printf("[%s] Finished", task)
			taskChan <- task
			return allTasksLocal
		}
	}
	return nil
}
*/

func Client(command *RemoteCommandClient, socketName string) int {
	log.Println("Client is running:", command.Args[0])
	var client net.Conn
	var err error
	// client, err = net.Dial("unix", socketName)
	client, err = net.Dial("tcp", socketName)
	if err != nil {
		log.Fatal(err)
	}
	transport, err := spdy.NewClientTransport(client)
	if err != nil {
		log.Fatal(err)
	}
	sender, err := transport.NewSendChannel()
	if err != nil {
		log.Fatal(err)
	}

	receiver, remoteSender := libchan.Pipe()
	//receiver, remoteSender := libchan.Pipe()
	command.StatusChan = remoteSender
	//receiver, _ := libchan.Pipe()

	err = sender.Send(command)
	if err != nil {
		log.Fatal(err)
	}

	response := &CommandResponse{}
	err = receiver.Receive(response)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("returning", response.Status)
	return response.Status

	return 0
}
