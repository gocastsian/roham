package grpc

import (
	"time"

	"google.golang.org/grpc"
)

type Config struct {
	Port               int           `koanf:"port"`
	NetworkType        string        `koanf:"type"`
	ShutDownCtxTimeout time.Duration `koanf:"shutdown_context_timeout"`
}

type RPCServer struct {
	Config Config
	Server *grpc.Server
}

func New(cfg Config) RPCServer {
	grpcServer := grpc.NewServer()

	return RPCServer{
		Server: grpcServer,
		Config: cfg,
	}
}

func (s RPCServer) Stop() {
	s.Server.GracefulStop()
}
