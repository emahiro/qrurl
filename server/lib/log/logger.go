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

// this is qrurl logger.
var logger *slog.Logger

var projectID = os.Getenv("GCP_PROJECT_ID")

type SpanIDKey struct{}
type TraceIDKey struct{}
type RequestTimeKey struct{}

func New() {
	logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))
}

type httpRequest struct {
	RequestMethod                  string   `json:"requestMethod,omitempty"`
	RequestUrl                     string   `json:"requestUrl,omitempty"`
	RequestSize                    string   `json:"requestSize,omitempty"`
	Status                         int      `json:"status,omitempty"`
	ResponseSize                   string   `json:"responseSize,omitempty"`
	UserAgent                      string   `json:"userAgent,omitempty"`
	RemoteIp                       string   `json:"remoteIp,omitempty"`
	ServerIp                       string   `json:"serverIp,omitempty"`
	Referer                        string   `json:"referer,omitempty"`
	Latency                        duration `json:"latency,omitempty"`
	CacheLookup                    bool     `json:"cacheLookup,omitempty"`
	CacheHit                       bool     `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer bool     `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 string   `json:"cacheFillBytes,omitempty"`
	Protocol                       string   `json:"protocol,omitempty"`
}

type duration struct {
	Nanos   int32 `json:"nanos,omitempty"`
	Seconds int64 `json:"seconds,omitempty"`
}

func makeDuration(d time.Duration) duration {
	nanos := d.Nanoseconds()
	secs := nanos / 1e9
	nanos -= secs * 1e9
	return duration{
		Nanos:   int32(nanos),
		Seconds: secs,
	}
}

func Requestf(ctx context.Context, rw http.ResponseWriter, r *http.Request) {

	spanID, ok := ctx.Value(SpanIDKey{}).(string)
	if !ok {
		spanID = ""
	}
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		traceID = ""
	}

	requestSize := fmt.Sprint(r.ContentLength)

	// response size を取得する
	var wb []byte
	n, err := rw.Write(wb)
	if err != nil {
		n = 0
	}
	responseSize := fmt.Sprint(n)

	requestTime, ok := ctx.Value(RequestTimeKey{}).(time.Time)
	if !ok {
		requestTime = time.Now()
	}
	duration := makeDuration(time.Since(requestTime))

	defer func() {
		logger.InfoCtx(ctx, "Default http request info",
			slog.String("logName", "projects/"+projectID+"/logs/qrurl-app%2FhttpRequestLog"),
			slog.String("severity", slog.LevelInfo.String()),
			slog.Any("httpRequest", httpRequest{
				RequestMethod: r.Method,
				RequestUrl:    r.URL.String(),
				RequestSize:   requestSize,
				ResponseSize:  responseSize,
				UserAgent:     r.UserAgent(),
				Protocol:      r.Proto,
				RemoteIp:      r.RemoteAddr,
				Referer:       r.Referer(),
				Latency:       duration,
			}),
			slog.Any("rawHttpHeader", r.Header),
			slog.Time("time", requestTime),
			slog.String("logging.googleapis.com/spanId", spanID),
			slog.String("logging.googleapis.com/trace", "projects/"+projectID+"/traces/"+traceID),
		)
	}()

}

type ConnectRequestInfo struct {
	Status      int
	RequestTime time.Time
	Duration    time.Duration
	Req         connect.AnyRequest
	Resp        connect.AnyResponse
}

func ConnectRequestf(ctx context.Context, info ConnectRequestInfo) {
	req := info.Req
	resp := info.Resp

	spanID, ok := ctx.Value(SpanIDKey{}).(string)
	if !ok {
		spanID = ""
	}
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		traceID = ""
	}

	requestSize, err := json.Marshal(req)
	if err != nil {
		requestSize = []byte{}
	}

	respSize, err := json.Marshal(resp)
	if err != nil {
		respSize = []byte{}
	}

	logger.InfoCtx(ctx, "Connect request info",
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app%2FconnectRequestLog"),
		slog.String("severity", slog.LevelInfo.String()),
		slog.Any("httpRequest", httpRequest{
			RequestMethod: req.HTTPMethod(),
			Status:        info.Status,
			RequestUrl:    "https://" + req.Header().Get("Host") + req.Spec().Procedure,
			RequestSize:   fmt.Sprint(len(requestSize)),
			UserAgent:     req.Header().Get("User-Agent"),
			Protocol:      req.Header().Get("Protocol"),
			RemoteIp:      req.Header().Get("X-Forwarded-For"),
			ServerIp:      req.Peer().Addr,
			ResponseSize:  fmt.Sprint(len(respSize)),
			Latency:       makeDuration(info.Duration),
		}),
		slog.Time("time", info.RequestTime),
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
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app%2FinfoLog"),
		slog.String("severity", slog.LevelInfo.String()),
		slog.Time("time", now),
		slog.String("message", msg),
		slog.String("logging.googleapis.com/spanId", spanID),
		slog.String("logging.googleapis.com/trace", "projects/"+projectID+"/traces/"+traceID),
	)
}
