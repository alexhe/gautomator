package flue

import (
	"io"
	"log"
	"net"
	"os/exec"
	"syscall"

	"github.com/docker/libchan"
	"github.com/docker/libchan/spdy"
)

// RemoteCommand is the run parameters to be executed remotely
type RemoteCommand struct {
	Cmd        string
	Args       []string
	Stdin      io.Writer
	Stdout     io.Reader
	Stderr     io.Reader
	StatusChan libchan.Sender
}

// CommandResponse is the returned response object from the remote execution
type CommandResponse struct {
	Status int
}

func Server(socketName string) {
	var listener net.Listener
	var err error
	listener, err = net.Listen("unix", socketName)
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
						command := &RemoteCommand{}
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

func Client(command RemoteCommand, socketName string) int {
	var client net.Conn
	var err error
	client, err = net.Dial("unix", socketName)
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

	err = sender.Send(command)
	if err != nil {
		log.Fatal(err)
	}

	response := &CommandResponse{}
	err = receiver.Receive(response)
	if err != nil {
		log.Fatal(err)
	}

	return response.Status
}
