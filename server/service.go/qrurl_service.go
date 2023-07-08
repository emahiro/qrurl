package service

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/bufbuild/connect-go"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
	"golang.org/x/exp/slog"

	qrurlv1 "github.com/emahiro/qrurl/server/gen/proto/qrurl/v1"
)

type QrUrlService struct{}

func (s *QrUrlService) PostCode(
	ctx context.Context,
	req *connect.Request[qrurlv1.PostCodeRequest],
) (resp *connect.Response[qrurlv1.PostCodeResponse], err error) {

	buf := bytes.NewBuffer(req.Msg.Image)
	img, _, err := image.Decode(buf)
	if err != nil {
		slog.ErrorCtx(ctx, "image.Decode", "err", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		slog.ErrorCtx(ctx, "NewBinaryBitmapFromImage", "err", err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	reader := qrcode.NewQRCodeReader()
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		slog.ErrorCtx(ctx, "reader.Decode", "err", err)
		return nil, connect.NewError(connect.CodeUnknown, err)
	}

	qrurlResp := &qrurlv1.PostCodeResponse{
		Url: result.GetText(),
	}

	resp = connect.NewResponse(qrurlResp)
	return resp, nil
}
