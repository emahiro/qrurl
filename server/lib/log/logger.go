package log

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/bufbuild/connect-go"
	"golang.org/x/exp/slog"
	"google.golang.org/protobuf/types/known/durationpb"
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
	RequestMethod                  string               `json:"requestMethod,omitempty"`
	RequestUrl                     string               `json:"requestUrl,omitempty"`
	RequestSize                    string               `json:"requestSize,omitempty"`
	Status                         int                  `json:"status,omitempty"`
	ResponseSize                   string               `json:"responseSize,omitempty"`
	UserAgent                      string               `json:"userAgent,omitempty"`
	RemoteIp                       string               `json:"remoteIp,omitempty"`
	ServerIp                       string               `json:"serverIp,omitempty"`
	Referer                        string               `json:"referer,omitempty"`
	Latency                        *durationpb.Duration `json:"latency,omitempty"`
	CacheLookup                    bool                 `json:"cacheLookup,omitempty"`
	CacheHit                       bool                 `json:"cacheHit,omitempty"`
	CacheValidatedWithOriginServer bool                 `json:"cacheValidatedWithOriginServer,omitempty"`
	CacheFillBytes                 string               `json:"cacheFillBytes,omitempty"`
	Protocol                       string               `json:"protocol,omitempty"`
}

func binarySize(v any) int {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(v); err != nil {
		return 0
	}
	return binary.Size(buf.Bytes())
}

type HTTPRequestLogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       uint64
}

func NewLogHTTPResponseWriter(rw http.ResponseWriter) *HTTPRequestLogResponseWriter {
	return &HTTPRequestLogResponseWriter{rw, http.StatusOK, 0}
}

func (lrw *HTTPRequestLogResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *HTTPRequestLogResponseWriter) Write(buf []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(buf)
	atomic.AddUint64(&lrw.size, uint64(n))
	return n, err
}

func defaultLogAttrs(severity slog.Level, traceID, spanID, message string) []slog.Attr {
	return []slog.Attr{
		slog.String("severity", severity.String()),
		slog.String("message", message),
		slog.String("logging.googleapis.com/spanId", spanID),
		slog.String("logging.googleapis.com/trace", "projects/"+projectID+"/traces/"+traceID),
	}
}

func Requestf(ctx context.Context, rw *HTTPRequestLogResponseWriter, r *http.Request) {

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
	if _, err := rw.Write([]byte{}); err != nil {
		rw.size = 0
	}
	responseSize := fmt.Sprint(rw.size)

	requestTime, ok := ctx.Value(RequestTimeKey{}).(time.Time)
	if !ok {
		requestTime = time.Now()
	}

	msg := "Default http request info"
	attrs := append(
		defaultLogAttrs(slog.LevelInfo, traceID, spanID, msg),
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app%2FhttpRequestLog"),
		slog.Time("time", requestTime),
		slog.Any("httpRequest", httpRequest{
			RequestMethod: r.Method,
			Status:        rw.statusCode,
			RequestUrl:    "https://" + r.Host + r.URL.String(),
			RequestSize:   requestSize,
			ResponseSize:  responseSize,
			UserAgent:     r.UserAgent(),
			Protocol:      r.Proto,
			RemoteIp:      r.RemoteAddr,
			Referer:       r.Referer(),
			Latency:       durationpb.New(time.Since(requestTime)),
		}),
		slog.Any("rawHttpHeader", r.Header), // ドキュメントに記載されてないフィールドは jsonPayload の内部に自動的に入る
	)
	logger.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

type ConnectRequestInfo struct {
	Status int
	Req    connect.AnyRequest
	Resp   connect.AnyResponse
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

	requestTime, ok := ctx.Value(RequestTimeKey{}).(time.Time)
	if !ok {
		requestTime = time.Now()
	}

	msg := "Connect request info"
	attrs := append(
		defaultLogAttrs(slog.LevelInfo, traceID, spanID, msg),
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app%2FconnectRequestLog"),
		slog.Any("httpRequest", httpRequest{
			RequestMethod: req.HTTPMethod(),
			Status:        info.Status,
			RequestUrl:    "https://" + req.Header().Get("Host") + req.Spec().Procedure,
			RequestSize:   fmt.Sprint(binarySize(req)),
			UserAgent:     req.Header().Get("User-Agent"),
			Protocol:      req.Header().Get("Protocol"),
			RemoteIp:      req.Header().Get("X-Forwarded-For"),
			ServerIp:      req.Peer().Addr,
			ResponseSize:  fmt.Sprint(binarySize(resp)),
			Latency:       durationpb.New(time.Since(requestTime)),
		}),
		slog.Time("time", requestTime),
		slog.Any("rawHttpHeader", req.Header()), // ドキュメントに記載されてないフィールドは jsonPayload の内部に自動的に入る
	)
	logger.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
}

func Infof(ctx context.Context, format string, args ...any) {
	spanID, ok := ctx.Value(SpanIDKey{}).(string)
	if !ok {
		spanID = ""
	}
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		traceID = ""
	}
	msg := fmt.Sprintf(format, args...)
	attrs := append(
		defaultLogAttrs(slog.LevelInfo, traceID, spanID, msg),
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app%2FErrorLog"),
		slog.Time("time", time.Now()),
	)
	logger.LogAttrs(ctx, slog.LevelInfo, "", attrs...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	spanID, ok := ctx.Value(SpanIDKey{}).(string)
	if !ok {
		spanID = ""
	}
	traceID, ok := ctx.Value(TraceIDKey{}).(string)
	if !ok {
		traceID = ""
	}
	msg := fmt.Sprintf(format, args...)
	attrs := append(
		defaultLogAttrs(slog.LevelError, traceID, spanID, msg),
		slog.String("logName", "projects/"+projectID+"/logs/qrurl-app%2FErrorLog"),
		slog.Time("time", time.Now()),
	)
	logger.LogAttrs(ctx, slog.LevelError, "", attrs...)

}
