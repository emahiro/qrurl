package service

import (
	"context"

	"github.com/bufbuild/connect-go"

	linewebhookv1 "github.com/emahiro/qrurl/server/gen/proto/webhook/v1"
)

type LineWebhookService struct{}

func (s *LineWebhookService) LineWebhook(ctx context.Context, req *connect.Request[linewebhookv1.LineWebhookRequest]) (*connect.Response[linewebhookv1.LineWebhookResponse], error) {
	resp := &linewebhookv1.LineWebhookResponse{
		Text: "Ok",
	}
	return connect.NewResponse(resp), nil
}
