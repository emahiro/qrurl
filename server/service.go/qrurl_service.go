package service

import (
	"context"

	"github.com/bufbuild/connect-go"

	qrurlv1 "github.com/emahiro/qrurl/server/gen/proto/qrurl/v1"
)

type QrUrlService struct{}

func (s *QrUrlService) PostCode(
	ctx context.Context,
	req *connect.Request[qrurlv1.PostCodeRequest],
) (resp *connect.Response[qrurlv1.PostCodeResponse], err error) {
	qrurlResp := &qrurlv1.PostCodeResponse{
		Url: "test",
	}
	resp = connect.NewResponse(qrurlResp)
	return resp, nil
}
