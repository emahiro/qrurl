package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/bufbuild/connect-go"

	"github.com/emahiro/qrurl/server/gen/proto/ping/v1/pingv1connect"
	"github.com/emahiro/qrurl/server/gen/proto/qrurl/v1/qrurlv1connect"
	"github.com/emahiro/qrurl/server/gen/proto/webhook/v1/webhookv1connect"
	"github.com/emahiro/qrurl/server/intercepter"
	"github.com/emahiro/qrurl/server/service"
)

const addr = ":8080"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	intercepters := connect.WithInterceptors(
		intercepter.NewRequestLogIntercepter(),
	)

	mux := http.NewServeMux()
	mux.Handle(qrurlv1connect.NewQrUrlServiceHandler(&service.QrUrlService{}, intercepters))
	mux.Handle(pingv1connect.NewPingServiceHandler(&service.PingService{}, intercepters))
	mux.Handle(webhookv1connect.NewLineWebhookServiceHandler(&service.LineWebhookService{}, intercepters))
	mux.HandleFunc("/v1/webhook/line", func(w http.ResponseWriter, r *http.Request) {
		v := map[string]any{}
		if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
			slog.ErrorCtx(ctx, "request body parse error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		slog.InfoCtx(r.Context(), "request body", "body", fmt.Sprintf("%+v", v))
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})
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
