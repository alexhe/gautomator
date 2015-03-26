package flue

import (
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
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

func RunTask(task string, allTasks chan *TopologyGraphStructure,wg *sync.WaitGroup) *TopologyGraphStructure {
	allTasksLocal := <-allTasks
	log.Println("Queuing:", task)
	if _, ok := allTasksLocal.waiter[task]; ok {
		log.Println("Waiting for the folowing tasks to finish:")

		for _, s := range allTasksLocal.waiter[task] {
			log.Println(s)
		}
	} else {
		log.Println("Let's go")
		command := &RemoteCommandClient{
			Cmd:    "sleep",
			Args:   []string{"5"},
			Stdin:  os.Stdin,
			Stdout: os.Stdout,
			Stderr: os.Stderr,
			//StatusChan: remoteSender,
		}
		//Client(command, "/tmp/mysocket.sock")
		Client(command, "localhost:5678")
		log.Println("Finished:", task)
	}
	    wg.Done()
	return allTasksLocal
}
func Client(command *RemoteCommandClient, socketName string) int {
	log.Println("Entering the client goroutine")
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
}
