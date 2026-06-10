PROTO_GEN_DIR := pkg/proto_gen
PROTO_GEN_NAME := proto_gen

.PHONY: run-broker run-broker-grpc run-publisher run-publisher-grpc \
        run-subscriber run-subscriber-grpc \
        build-broker build-publisher build-subscriber build-all \
        proto-gen test

run-broker:
	go run cmd/broker/rpc/main.go

run-broker-grpc:
	go run cmd/broker/grpc/main.go

run-publisher:
	go run cmd/publisher/rpc/main.go

run-publisher-grpc:
	go run cmd/publisher/grpc/main.go

run-subscriber:
	go run cmd/subscriber/rpc/main.go

run-subscriber-grpc:
	go run cmd/subscriber/grpc/main.go

build-broker:
	go build -o bin/broker cmd/broker/rpc/main.go
	go build -o bin/broker_grpc cmd/broker/grpc/main.go

build-publisher:
	go build -o bin/publisher cmd/publisher/rpc/main.go
	go build -o bin/publisher_grpc cmd/publisher/grpc/main.go

build-subscriber:
	go build -o bin/subscriber cmd/subscriber/rpc/main.go
	go build -o bin/subscriber_grpc cmd/subscriber/grpc/main.go

build-all: build-broker build-publisher build-subscriber

test:
	go test -timeout 30s -race ./...

proto-gen:
	mkdir -p $(PROTO_GEN_DIR)
	protoc --go_out=$(PROTO_GEN_DIR) --go-grpc_out=$(PROTO_GEN_DIR) protos/broker.proto
