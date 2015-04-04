# flue
A distributed agentless configuration and management tool written in go

This is my first project in Go.

What **flue** should be someday:
- a duct for binaries and configuration for deploying an infrastructure and/or a distributed application.
- a distributed configuration and management software.

*flue* aim to replicate and spread itself like a *flu*.

# Why go ?

- a compiled go project is a single binary, which make is good for an heterogenous infrastructure. You just copy the *flue* and it does not need anything else to run (no python dependency, no perl, no ruby, no whatever lib...)
- **Concurency** of course !
- It looks fast
- Not too difficult to learn (I will tell you that in a few weeks)
- A lot of library exists in the standard implementation

# WIP

## What's working:
- The task parsing (from a dot file)
- The remote task execution implementing the libchan from docker
- the web server and the REST api is in progress but working

## How to:
1) launch a server in a terminal: `go run agent/main.go`
2) launch a "client" from another terminal: `go run agent/main.go -dot=example.dot`
3) follow the stream on `http://localhost:8080/sigma.html`

## What's next:
- pass several tasks on the command line
- merge the tasks
- clean the web server and the REST code
- clean all the code
- implement the "node" concept

# How to

* grab a Go compiler from https://golang.org
* `go get github.com/owulveryck/flue`

# tests and developement

This work is developed on MacOS and FreeBSD.
It will be tested as well on Linux.

I will try to stick to pure go implementation to remain portable.

# Note
This project is in developpement.
It is a work in progress and a Feature Driven Developement.
On top of that, I'm learning go (and the developement). 

Be tolerant !
