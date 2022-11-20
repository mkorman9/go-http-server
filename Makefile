OUTPUT ?= go-http-server

.DEFAULT_GOAL := compile

compile:
	CGO_ENABLED=0 go build -o $(OUTPUT)

package:
	VERSION=$(shell make version) bash build/package.sh

version:
	@echo "1.0.0"

all: compile package
