package models

import (
	"fmt"
	"net/rpc"
)

type Transport string

const (
	DefaultTransport   Transport = "RPC"
	GrpcTransport      Transport = "GRPC"
	WebSocketTransport Transport = "WebSocket"
	RedisTransport     Transport = "Redis"
)

// BrokerConfig contains information for setting up a broker server
type BrokerConfig struct {
	Host            string
	Port            string
	Storage         StorageType
	StorageLocation *string
	Transport       Transport
}

func ConnectToRpcServer(host, port string) (*rpc.Client, error) {
	return rpc.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
}
