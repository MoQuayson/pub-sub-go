package server

import "net/rpc"

type Server interface {
	Register(service any) error
	RegisterName(name string, service any)
	Call(serviceMethod string, args, reply any)
}

type RpcServer struct {
	Server
}

func NewRpcServer() Server {
	srv := RpcServer{}

	return &srv
}

func (s *RpcServer) Register(service any) error {
	if err := rpc.Register(service); err != nil {
		return err
	}

	return nil
}
