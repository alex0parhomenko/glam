-include .env

VERSION := $(shell git rev-parse --short HEAD)
PROJECTNAME := $(shell basename "$(PWD)")

# Go related variables.
GOBASE := $(shell pwd)
GOBIN := $(GOBASE)/bin
GOFILES := cmd/main.go

# Use linker flags to provide version/build
LDFLAGS=-ldflags "-X=main.Version=$(VERSION)"

# Redirect error output to a file, so we can show it in development mode.
STDERR := /tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID := /tmp/.$(app).pid

# Make is verbose in Linux. Make it silent.
MAKEFLAGS += --silent

go-build:
	@echo "  >  Building binary $(app)..."
	go build $(LDFLAGS) -o $(GOBIN)/glam $(GOFILES)

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
