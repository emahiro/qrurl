package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/emahiro/qrurl/server/gen/proto/ping/v1/pingv1connect"
	"github.com/emahiro/qrurl/server/gen/proto/qrurl/v1/qrurlv1connect"
	"github.com/emahiro/qrurl/server/handler"
	"github.com/emahiro/qrurl/server/infra/firestore"
	"github.com/emahiro/qrurl/server/intercepter"
	"github.com/emahiro/qrurl/server/lib/line"
	"github.com/emahiro/qrurl/server/lib/log"
	"github.com/emahiro/qrurl/server/middleware"
	"github.com/emahiro/qrurl/server/service"
)

const addr = ":8080"

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// init logger
	log.New()

	// init line
	if err := line.NewBot(ctx, true); err != nil {
		log.Errorf(ctx, "failed to init line bot. err: %+v", err)
		panic(err)
	}
	// init firestore
	if err := firestore.New(ctx); err != nil {
		log.Errorf(ctx, "failed to init firestore. err: %+v", err)
		panic(err)
	}

	// cors setting
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedHeaders: []string{"*"},
	})

	mws := []func(http.Handler) http.Handler{
		middleware.SetTrace,
		middleware.RequestLog,
	}
	mux := http.NewServeMux()
	mux.Handle("/v1/webhook/line", middleware.Chain(
		http.HandlerFunc(handler.LineWebHookHandler),
		append(mws, middleware.VerifyChannelAccessToken, middleware.VerifyLine)...,
	))
	mux.Handle("/ping", middleware.Chain(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { fmt.Fprintf(w, "{\"message\": \"pong\"}") }), mws...))

	intercepters := connect.WithInterceptors(
		intercepter.NewRequestLogIntercepter(),
	)
	mux.Handle(qrurlv1connect.NewQrUrlServiceHandler(&service.QrUrlService{}, intercepters))
	mux.Handle(pingv1connect.NewPingServiceHandler(&service.PingService{}, intercepters))

	server := &http.Server{
		Addr:    addr,
		Handler: c.Handler(h2c.NewHandler(mux, &http2.Server{})),
	}
	go func() {
		<-ctx.Done()
		log.Infof(context.Background(), "shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			log.Errorf(context.Background(), "failed to shutdown server: %v", err)
			panic(err)
		}
	}()

	log.Infof(ctx, "server start! port: %v", "localhost"+addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		panic(err)
	}
}
