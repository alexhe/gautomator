# gautomator
A distributed agentless configuration and management tool written in go base on graph

This is my first project in Go.

What **gautomator** should be someday:
- a duct for binaries and configuration for deploying an infrastructure and/or a distributed application.
- a distributed configuration and management software.

*gautomator* aim to replicate and spread itself like a *flu*.

# Why go ?

- a compiled go project is a single binary, which make is good for an heterogenous infrastructure. You just copy the *gautomator* and it does not need anything else to run (no python dependency, no perl, no ruby, no whatever lib...)
- **Concurency** of course !
- It looks fast
- Not too difficult to learn (I will tell you that in a few weeks)
- A lot of library exists in the standard implementation

# WIP
gautomator is  in heavy developement, see the wiki and the issue for more information.
And for even more information about the implementation:

* twitter: @owulveryck

# How to

* grab a Go compiler from https://golang.org
* `go get github.com/owulveryck/gautomator`

# tests and developement

This work is developed on MacOS and FreeBSD.
It will be tested as well on Linux.

I will try to stick to pure go implementation to remain portable.

# Dependencies

I would stick to pure go implementation in order to get a single binary that may run "anywhere"

Go dependency so far:

* `github.com/docker/libchan` for the client/server communication
* `github.com/docker/libchan/spdy`
* `github.com/gonum/matrix/mat64` for the matrix manipulation (for the graph theory)
* `github.com/awalterschulze/gographviz` for parsing the topology files
* `github.com/gorilla/mux` for the webserver implementation


# Note
This project is in developpement.
It is a work in progress and a Feature Driven Developement.
On top of that, I'm learning go (and the developement). 

Be tolerant !
