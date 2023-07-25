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

var projectID = os.Getenv("GCP_PROJECT_ID")

type SpanIDKey struct{}
type TraceIDKey struct{}

func New() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

type httpRequest struct {
	RequestMethod                  string `json:"requestMethod,omitempty"`
	RequestUrl                     string `json:"requestUrl,omitempty"`
	RequestSize                    string `json:"requestSize,omitempty"`
	Status                         int    `json:"status,omitempty"`
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

	spanID, ok := ctx.Value(SpanIDKey{}).(string)
	if !ok {
		spanID = ""
	}
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		traceID = ""
	}

	logger.InfoCtx(ctx, "Default http request info",
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app-http-request"),
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
		slog.Any("rawHttpHeader", r.Header),
		slog.Time("time", now),
		slog.String("logging.googleapis.com/spanId", spanID),
		slog.String("logging.googleapis.com/trace", "projects/"+projectID+"/traces/"+traceID),
	)
}

func ConnectRequestf(ctx context.Context, status int, req connect.AnyRequest, resp connect.AnyResponse) {
	now := time.Now()

	spanID, ok := ctx.Value(SpanIDKey{}).(string)
	if !ok {
		spanID = ""
	}
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		traceID = ""
	}

	respSize, err := json.Marshal(resp)
	if err != nil {
		respSize = []byte{}
	}
	logger.InfoCtx(ctx, "Connect request info",
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app-connect-request"),
		slog.String("severity", slog.LevelInfo.String()),
		slog.Any("httpRequest", httpRequest{
			RequestMethod: req.HTTPMethod(),
			Status:        status,
			RequestUrl:    req.Header().Get("Host") + req.Spec().Procedure,
			RequestSize:   req.Header().Get("Content-Length"),
			UserAgent:     req.Header().Get("User-Agent"),
			Protocol:      req.Header().Get("Protocol"),
			RemoteIp:      req.Header().Get("X-Forwarded-For"),
			ServerIp:      req.Peer().Addr,
			ResponseSize:  fmt.Sprint(len(respSize)),
		}),
		slog.Time("time", now),
		slog.Any("rawHttpHeader", req.Header()), // ドキュメントに記載されてないフィールドは jsonPayload の内部に自動的に入る
		slog.String("logging.googleapis.com/spanId", spanID),
		slog.String("logging.googleapis.com/trace", "projects/"+projectID+"/traces/"+traceID),
	)
}

func Infof(ctx context.Context, format string, args ...any) {
	now := time.Now()

	spanID, ok := ctx.Value(SpanIDKey{}).(string)
	if !ok {
		spanID = ""
	}
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		traceID = ""
	}

	msg := fmt.Sprintf(format, args...)
	logger.LogAttrs(ctx, slog.LevelInfo, msg,
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app"),
		slog.String("severity", slog.LevelInfo.String()),
		slog.Time("time", now),
		slog.String("message", msg),
		slog.String("logging.googleapis.com/spanId", spanID),
		slog.String("logging.googleapis.com/trace", "projects/"+projectID+"/traces/"+traceID),
	)
}
