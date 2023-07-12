package service

import (
	"context"

	webhookv1 "github.com/emahiro/qrurl/server/gen/proto/webhook/v1"
)

type LineWebhookService struct{}

func (s *LineWebhookService) Line(ctx context.Context, req *webhookv1.LineWebhookRequest) (*webhookv1.LineWebhookResponse, error) {
	resp := &webhookv1.LineWebhookResponse{
		Text: "Ok",
	}
	return resp, nil
}
