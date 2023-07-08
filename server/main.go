package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"golang.org/x/exp/slog"

	qrurlv1 "github.com/emahiro/qrurl/server/gen/proto/qrurl/v1"
	"github.com/emahiro/qrurl/server/gen/proto/qrurl/v1/qrurlv1connect"
)

const addr = ":8080"

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

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.Handle(qrurlv1connect.NewQrUrlServiceHandler(&QrUrlService{}))
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	go func() {
		<-ctx.Done()
		if err := server.Close(); err != nil {
			panic(err)
		}
	}()

	slog.InfoCtx(ctx, "server start", "port", fmt.Sprintf("localhost%s", addr))
	if err := server.ListenAndServe(); err != nil {
		panic(err)
	}
}
