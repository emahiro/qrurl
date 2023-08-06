package service

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"

	pingv1 "github.com/emahiro/qrurl/server/gen/proto/ping/v1"
	"github.com/emahiro/qrurl/server/lib/log"
)

type PingService struct{}

func (s *PingService) Ping(ctx context.Context, req *connect.Request[pingv1.PingRequest]) (*connect.Response[pingv1.PingResponse], error) {
	pingResp := &pingv1.PingResponse{
		Message: "pong",
	}
	err := log.WithStackTracef(errors.New("test"), "test")
	log.Infof(ctx, "err: %+v", err)
	resp := connect.NewResponse(pingResp)
	return resp, nil
}
