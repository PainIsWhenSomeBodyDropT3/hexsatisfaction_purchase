package grpc

import (
	"github.com/JesusG2000/hexsatisfaction/pkg/grpc/api"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

type Checker struct {
	api.ExistanceClient
}

func NewGRPCClient(addr string) (*Checker, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "couldn't dial connection with gprc")
	}

	return &Checker{api.NewExistanceClient(conn)}, nil
}
