package lib

import (
	"bytes"
	"context"
	"image"
	"log/slog"

	"github.com/bufbuild/connect-go"
	"github.com/makiuchi-d/gozxing"
	"github.com/makiuchi-d/gozxing/qrcode"
)

func DecodeQrCode(ctx context.Context, b []byte) (string, error) {
	buf := bytes.NewBuffer(b)
	img, _, err := image.Decode(buf)
	if err != nil {
		slog.ErrorCtx(ctx, "image.Decode", "err", err)
		return "", connect.NewError(connect.CodeInvalidArgument, err)
	}

	bmp, err := gozxing.NewBinaryBitmapFromImage(img)
	if err != nil {
		slog.ErrorCtx(ctx, "NewBinaryBitmapFromImage", "err", err)
		return "", connect.NewError(connect.CodeUnknown, err)
	}

	reader := qrcode.NewQRCodeReader()
	result, err := reader.Decode(bmp, nil)
	if err != nil {
		slog.ErrorCtx(ctx, "reader.Decode", "err", err)
		return "", connect.NewError(connect.CodeUnknown, err)
	}
	return result.GetText(), nil
}
