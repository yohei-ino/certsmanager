# Makefile

# This Makefile is used to build the project and manage dependencies.


# ==============================================================================
# Makefile configuration
# ==============================================================================

DBG_MAKEFILE ?=
ifeq ($(DBG_MAKEFILE),1)
    $(warning ***** starting Makefile for goal(s) "$(MAKECMDGOALS)")
    $(warning ***** $(shell date))
else
    # デバッグモードでない場合はサイレントモードにする
    MAKEFLAGS += -s
endif

PROJECT_NAME := certsmanager
PROJECT_VERSION := 1.0.0
PROJECT_DESCRIPTION := A simple certificate manager written in Go.
PROJECT_AUTHOR := yohei-ino
PROJECT_LICENSE := MIT
PROJECT_ROOT := .
BIN_DIR := $(PROJECT_ROOT)/bin
SRC_DIR := $(PROJECT_ROOT)/src


.PHONY: all init build test clean

all: build

init:
	@echo "Initializing Go module..."
	go mod init $(PROJECT_NAME)
	go mod tidy
	@echo "Initialization completed!"

build: main

main: $(SRC_DIR)/main.go
	mkdir -p $(BIN_DIR)
	@echo "Building $(PROJECT_NAME)..."
	go build -o $(BIN_DIR)/$(PROJECT_NAME) $(SRC_DIR)/main.go
	@echo "Build completed!"

test:
	@echo "Running tests..."
	go test -v ./...
	@echo "Tests completed!"

clean:
	rm -rf $(BIN_DIR)
	@echo "Cleaned up."

