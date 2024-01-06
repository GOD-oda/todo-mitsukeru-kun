#!/bin/bash

NAME="todo-mitsukeru-kun"

env GOOS=linux GOARCH=amd64 go build -o bin/dist/linux/${NAME}
env GOOS=darwin GOARCH=amd64 go build -o bin/dist/darwin/${NAME}