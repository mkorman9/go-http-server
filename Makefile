OUTPUT ?= go-http-server

.DEFAULT_GOAL := compile

compile:
	CGO_ENABLED=0 go build -o $(OUTPUT)

all: compile
