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
	now := time.Now()
	header, err := json.Marshal(r.Header)
	if err != nil {
		slog.ErrorCtx(ctx, "failed to marshal header", "err", err)
		header = []byte(`{}`)
	}
	hr, err := json.Marshal(httpRequest{
		RequestMethod: r.Method,
		RequestUrl:    r.URL.String(),
		RequestSize:   fmt.Sprintf("%d", r.ContentLength),
		UserAgent:     r.UserAgent(),
		Protocol:      r.Proto,
		RemoteIp:      r.RemoteAddr,
		Referer:       r.Referer(),
	})
	if err != nil {
		slog.ErrorCtx(ctx, "failed to marshal httpRequest", "err", err)
		hr = []byte(`{}`)
	}
	logger.InfoCtx(ctx, "Default http request info",
		slog.String("severity", slog.LevelInfo.String()),
		slog.String("httpRequest", string(hr)),
		slog.Any("jsonPayload", header),
		slog.Time("time", now),
	)
}

func ConnectRequestf(ctx context.Context, r connect.AnyRequest) {
	now := time.Now()
	hr, err := json.Marshal(httpRequest{
		RequestMethod: r.HTTPMethod(),
		RequestUrl:    r.Spec().Procedure,
		RequestSize:   r.Header().Get("Content-Length"),
		UserAgent:     r.Header().Get("User-Agent"),
		Protocol:      r.Peer().Protocol,
		RemoteIp:      r.Peer().Addr,
	})
	if err != nil {
		slog.ErrorCtx(ctx, "failed to marshal httpRequest", "err", err)
		hr = []byte(`{}`)
	}
	logger.InfoCtx(ctx, "Connect request info",
		slog.String("severity", slog.LevelInfo.String()),
		slog.String("httpRequest", string(hr)),
		slog.Any("jsonPayload", slog.AnyValue(r.Header())),
		slog.Time("time", now),
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
		slog.String("textPayload", msg),
		slog.String("logging.googleapis.com/insertId", "0"),
		slog.String("logging.googleapis.com/spanId", "0"),
		slog.String("logging.googleapis.com/trace", "0"),
	)
}
