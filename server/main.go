package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/go-chi/chi/v5"
	"golang.org/x/exp/slog"

	"github.com/emahiro/qrurl/server/gen/proto/ping/v1/pingv1connect"
	"github.com/emahiro/qrurl/server/gen/proto/qrurl/v1/qrurlv1connect"
	"github.com/emahiro/qrurl/server/intercepter"
	"github.com/emahiro/qrurl/server/service"
)

const addr = ":8080"

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	router := chi.NewRouter()
	// 標準の http handler
	router.Group(func(r chi.Router) {
		r.HandleFunc("/v1/webhook/line", func(w http.ResponseWriter, r *http.Request) {
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
