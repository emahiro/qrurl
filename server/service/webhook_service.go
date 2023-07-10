package service

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"golang.org/x/exp/slog"

	linewebhookv1 "github.com/emahiro/qrurl/server/gen/proto/webhook/v1"
)

type LineWebhookService struct{}

func (s *LineWebhookService) LineWebhook(ctx context.Context, req *connect.Request[linewebhookv1.LineWebhookRequest]) (*connect.Response[linewebhookv1.LineWebhookResponse], error) {
	body := req.Msg
	slog.InfoCtx(ctx, "request", fmt.Sprintf("%v", body))
	resp := &linewebhookv1.LineWebhookResponse{
		Text: "Ok",
	}
	return connect.NewResponse(resp), nil
}
