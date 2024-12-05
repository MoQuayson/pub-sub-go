PROTO_GEN_DIR=pkg/proto_gen
PROTO_GEN_NAME=proto_gen

run-broker:
	go run cmd/broker/rpc/main.go

build-broker:
	go build -o bin/broker.exe cmd/broker/rpc/main.go
	go build -o bin/broker_grpc.exe cmd/broker/grpc/main.go

run-publisher:
	go run cmd/publisher/main.go

build-publisher:
	go build -o bin/publisher.exe cmd/publisher/rpc/main.go
	go build -o bin/publisher_grpc.exe cmd/publisher/grpc/main.go

run-subscriber:
	go run cmd/subscriber/main.go

build-subscriber:
	go build -o bin/subscriber.exe cmd/subscriber/rpc/main.go
	go build -o bin/subscriber_grpc.exe cmd/subscriber/grpc/main.go

build-all: build-broker build-publisher build-subscriber

proto-gen:
	cd pkg && mkdir ${PROTO_GEN_NAME}
	protoc --go_out=${PROTO_GEN_DIR} --go-grpc_out=${PROTO_GEN_DIR} protos/broker.proto

