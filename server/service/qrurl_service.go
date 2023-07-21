package service

import (
	"context"
	_ "image/jpeg"
	_ "image/png"

	"github.com/bufbuild/connect-go"

	qrurlv1 "github.com/emahiro/qrurl/server/gen/proto/qrurl/v1"
	"github.com/emahiro/qrurl/server/lib"
)

type QrUrlService struct{}

func (s *QrUrlService) PostQrCode(
	ctx context.Context,
	req *connect.Request[qrurlv1.PostQrCodeRequest],
) (resp *connect.Response[qrurlv1.PostQrCodeResponse], err error) {
	url, err := lib.DecodeQrCode(ctx, req.Msg.Image)
	if err != nil {
		return nil, err
	}
	qrurlResp := &qrurlv1.PostQrCodeResponse{
		Url: url,
	}
	resp = connect.NewResponse(qrurlResp)
	return resp, nil
}
