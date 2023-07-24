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
	RequestMethod                  string `json:"requestMethod,omitempty"`
	RequestUrl                     string `json:"requestUrl,omitempty"`
	RequestSize                    string `json:"requestSize,omitempty"`
	Status                         int64  `json:"status,omitempty"`
	ResponseSize                   string `json:"responseSize,omitempty"`
	UserAgent                      string `json:"userAgent,omitempty"`
	RemoteIp                       string `json:"remoteIp,omitempty"`
	ServerIp                       string `json:"serverIp,omitempty"`
	Referer                        string `json:"referer,omitempty"`
	Latency                        string `json:"latency,omitempty"`
	CacheLookup                    bool   `json:"cacheLookup,omitempty"`
	CacheHit                       bool   `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer bool   `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 string `json:"cacheFillBytes,omitempty"`
	Protocol                       string `json:"protocol,omitempty"`
}

func Requestf(ctx context.Context, r *http.Request) {
	now := time.Now()
	header, err := json.Marshal(r.Header)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to marshal header", "err", err)
		header = []byte(`{}`)
	}
	logger.InfoCtx(ctx, "Default http request info",
		slog.String("severity", slog.LevelInfo.String()),
		slog.Any("httpRequest", httpRequest{
			RequestMethod: r.Method,
			RequestUrl:    r.URL.String(),
			RequestSize:   fmt.Sprintf("%d", r.ContentLength),
			UserAgent:     r.UserAgent(),
			Protocol:      r.Proto,
			RemoteIp:      r.RemoteAddr,
			Referer:       r.Referer(),
		}),
		slog.Any("jsonPayload", map[string]any{
			"header": header,
		}),
		slog.Time("time", now),
	)
}

func ConnectRequestf(ctx context.Context, r connect.AnyRequest) {
	now := time.Now()
	logger.InfoCtx(ctx, "Connect request info",
		slog.String("severity", slog.LevelInfo.String()),
		slog.Any("httpRequest", httpRequest{
			RequestMethod: r.HTTPMethod(),
			RequestUrl:    r.Spec().Procedure,
			RequestSize:   r.Header().Get("Content-Length"),
			UserAgent:     r.Header().Get("User-Agent"),
			Protocol:      r.Peer().Protocol,
			RemoteIp:      r.Peer().Addr,
		}),
		slog.Time("time", now),
		slog.Any("jsonPayload", map[string]any{
			"httpHeader": r.Header(),
		}),
		slog.String("logging.googleapis.com/insertId", "0"),
		slog.String("logging.googleapis.com/spanId", "0"),
		slog.String("logging.googleapis.com/trace", "0"),
	)
}

func Infof(ctx context.Context, format string, args ...any) {
	now := time.Now()
	msg := fmt.Sprintf(format, args...)
	logger.LogAttrs(ctx, slog.LevelInfo, msg,
		slog.String("severity", slog.LevelInfo.String()),
		slog.Time("time", now),
		slog.String("logging.googleapis.com/insertId", "0"),
		slog.String("logging.googleapis.com/spanId", "0"),
		slog.String("logging.googleapis.com/trace", "0"),
	)
}
