package src

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"learning/user_agent/pb"
	"time"
)

const ContextReqUUid = "req_uuid"

type NewMiddlewareServer func(service Service) Service

type logMiddlewareServer struct {
	logger *zap.Logger
	next   Service
}

func NewLogMiddlewareServer(log *zap.Logger) NewMiddlewareServer {
	return func(service Service) Service {
		return logMiddlewareServer{
			log,
			service,
		}
	}
}

func (l logMiddlewareServer) Login(ctx context.Context, in *pb.Login) (out *pb.LoginAck, err error) {
	defer func(start time.Time) {
		l.logger.Debug(fmt.Sprint(ctx.Value(ContextReqUUid)), zap.Any("调用 Login logMiddlewareServer", "Login"), zap.Any("req", in), zap.Any("res", out), zap.Any("time", time.Since(start)), zap.Any("err", err))
	}(time.Now())
	out, err = l.next.Login(ctx, in)
	return
}
