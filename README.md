# flue
A piece of code in go

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

# Principles

I separate the topology in layers.
By now, I count 3 layers, but there may be more in the future.
- Layer 1: the infrastructure layer. Here we are talking about os, ssh, users, ports and so on.
- Layer 2: the product layer. Here we are talking about apache, nginx, weblogic, jboss, ...
- Layer 3: the middleware layer. Here we are talking about applicative architecture: producer, consumer, database, webserver
- Layer 7: the applicative layer. Here we are talking about ear, war, zip, html, css, ...

**Flue** should be able to deploy any layer, and any composant of the layer.
