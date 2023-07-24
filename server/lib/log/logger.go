package log

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/bufbuild/connect-go"
	"golang.org/x/exp/slog"
)

var logger *slog.Logger

func New() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

type httpRequest struct {
	RequestMethod                  string `json:"requestMethod"`
	RequestUrl                     string `json:"requestUrl"`
	RequestSize                    string `json:"requestSize"`
	Status                         int64  `json:"status"`
	ResponseSize                   string `json:"responseSize"`
	UserAgent                      string `json:"userAgent"`
	RemoteIp                       string `json:"remoteIp"`
	ServerIp                       string `json:"serverIp"`
	Referer                        string `json:"referer"`
	Latency                        string `json:"latency"`
	CacheLookup                    bool   `json:"cacheLookup"`
	CacheHit                       bool   `json:"cacheHit"`
	CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer"`
	CacheFillBytes                 string `json:"cacheFillBytes"`
	Protocol                       string `json:"protocol"`
}

func Requestf(ctx context.Context, r *http.Request) {
	header, err := json.Marshal(r.Header)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to marshal header", "err", err)
		header = []byte(`{}`)
	}
	logger.InfoCtx(ctx, "http request info",
		"severity", slog.LevelInfo,
		"httpRequest", httpRequest{
			RequestMethod: r.Method,
			RequestUrl:    r.URL.String(),
			RequestSize:   fmt.Sprintf("%d", r.ContentLength),
			UserAgent:     r.UserAgent(),
			Protocol:      r.Proto,
			RemoteIp:      r.RemoteAddr,
			Referer:       r.Referer(),
		},
		"jsonPayload", header,
	)
}

func ConnectRequestf(ctx context.Context, r connect.AnyRequest) {
	logger.InfoCtx(ctx, "this is connect request info",
		"severity", slog.LevelInfo,
		"httpRequest", httpRequest{
			RequestMethod: r.HTTPMethod(),
			RequestUrl:    r.Spec().Procedure,
			RequestSize:   r.Header().Get("Content-Length"),
			UserAgent:     r.Header().Get("User-Agent"),
			Protocol:      r.Peer().Protocol,
			RemoteIp:      r.Peer().Addr,
		},
		"jsonPayload", r.Header(),
	)
}

func Infof(ctx context.Context, format string, args ...any) {
	now := time.Now()
	msg := fmt.Sprintf(format, args...)
	logger.InfoCtx(ctx, msg,
		"severity", slog.LevelInfo,
		"time", now.Format(time.RFC3339Nano),
		"message", msg,
		"insertId", "0",
		"spanId", "0",
		"trace", "0",
	)
}
