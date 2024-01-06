FROM golang:1.21.1-bullseye AS dev

ENV ROOT /todo-mitsukeru-kun
ENV CGO_ENABLED 0

WORKDIR ${ROOT}
RUN apt update  \
    && apt install git

COPY src/go.mod ./
RUN go mod download
COPY src .
