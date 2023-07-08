package main

import (
	"context"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/emahiro/qrurl/server/gen/proto/qrurl/v1/qrurlv1connect"
	"github.com/emahiro/qrurl/server/service.go"
)

const addr = ":8080"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := http.NewServeMux()
	mux.Handle(qrurlv1connect.NewQrUrlServiceHandler(&service.QrUrlService{}))
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
