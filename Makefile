BINARY_NAME=BST

all: 
	go build -o $(BINARY_NAME) src/*.go

build:
	go build -o $(BINARY_NAME) src/*.go