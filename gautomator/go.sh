#!/bin/sh

GOPATH=$PWD/../../../../../
dot -Tpng example.dot > ../htdocs/static/example.png
go run main.go example.dot
