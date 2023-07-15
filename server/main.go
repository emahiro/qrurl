package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/emahiro/qrurl/server/gen/proto/ping/v1/pingv1connect"
	"github.com/emahiro/qrurl/server/gen/proto/qrurl/v1/qrurlv1connect"
	"github.com/emahiro/qrurl/server/handler"
	"github.com/emahiro/qrurl/server/infra/firestore"
	"github.com/emahiro/qrurl/server/intercepter"
	"github.com/emahiro/qrurl/server/lib/line"
	"github.com/emahiro/qrurl/server/service"
)

const addr = ":8080"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// init line
	if err := line.NewBot(ctx, true); err != nil {
		slog.ErrorCtx(ctx, "failed to init line bot", "err", err)
		panic(err)
	}
	// init firestore
	if err := firestore.New(ctx); err != nil {
		slog.ErrorCtx(ctx, "failed to init firestore", "err", err)
		panic(err)
	}

	router := chi.NewRouter()
	// 標準の http handler
	router.Group(func(r chi.Router) {
		r.Use(
			intercepter.VerifyLine(),
			intercepter.VerifyChannelAccessToken(),
		)
		r.HandleFunc("/v1/webhook/line", handler.LineWebHookHandler)
	})

	// for GRPC
	router.Group(func(r chi.Router) {
		intercepters := connect.WithInterceptors(
			intercepter.NewRequestLogIntercepter(),
		)
		r.Handle(qrurlv1connect.NewQrUrlServiceHandler(&service.QrUrlService{}, intercepters))
		r.Handle(pingv1connect.NewPingServiceHandler(&service.PingService{}, intercepters))
	})

	server := &http.Server{
		Addr:    addr,
		Handler: router,
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
