package service

import (
	"context"
	"encoding/base64"
	_ "image/jpeg"
	_ "image/png"

	"github.com/bufbuild/connect-go"

	qrurlv1 "github.com/emahiro/qrurl/server/gen/proto/qrurl/v1"
	"github.com/emahiro/qrurl/server/lib"
	"github.com/emahiro/qrurl/server/lib/log"
)

type QrUrlService struct{}

func (s *QrUrlService) PostQrCode(
	ctx context.Context,
	req *connect.Request[qrurlv1.PostQrCodeRequest],
) (resp *connect.Response[qrurlv1.PostQrCodeResponse], err error) {
	b, err := base64.StdEncoding.DecodeString(req.Msg.Image)
	if err != nil {
		return nil, log.WithStackTracef(err, "decode error. image: %v", req.Msg.Image)
	}
	url, err := lib.DecodeQrCode(ctx, b)
	if err != nil {
		return nil, log.WithStackTracef(err, "decode error. image: %v", req.Msg.Image)
	}
	qrurlResp := &qrurlv1.PostQrCodeResponse{
		Url: url,
	}
	resp = connect.NewResponse(qrurlResp)
	return resp, nil
}
