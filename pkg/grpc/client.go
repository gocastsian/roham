package grpc

import (
	"fmt"
	"google.golang.org/grpc"
)

type Client struct {
	Host string `koanf:"host"`
	Port int    `koanf:"port"`
}

func NewClient(cfg Client) (*grpc.ClientConn, error) {
	address := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.NewClient(address)

	fmt.Println(err)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return conn, nil
}
