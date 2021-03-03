package src

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"learning/user_agent/pb"
	"learning/utils"
)

type Service interface {
	Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error)
}

type baseServer struct {
	logger *zap.Logger
}

func NewService(log *zap.Logger) Service {
	var server Service
	server = &baseServer{log}
	server = NewLogMiddlewareServer(log)(server)
	return server
}

func (s baseServer) Login(ctx context.Context, in *pb.Login) (ack *pb.LoginAck, err error) {
	if in.Account != "shier" || in.Password != "123123" {
		err = errors.New("用户信息错误")
		return
	}
	ack = &pb.LoginAck{}
	ack.Token, err = utils.CreateJwtToken(in.Account, 1)
	return
}
