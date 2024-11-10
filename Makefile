run-broker:
	go run cmd/broker/main.go

build-broker:
	go build -o bin/broker.exe cmd/broker/main.go

run-publisher:
	go run cmd/publisher/main.go

build-publisher:
	go build -o bin/publisher.exe cmd/publisher/main.go

run-subscriber:
	go run cmd/subscriber/main.go

build-subscriber:
	go build -o bin/subscriber.exe cmd/subscriber/main.go

build-all: build-broker build-publisher build-subscriber
