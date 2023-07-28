package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/emahiro/qrurl/server/lib/log"
)

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestTime := time.Now()
		ctx = context.WithValue(ctx, log.RequestTimeKey{}, requestTime)

		xcTraceCtx := r.Header.Get("X-Cloud-Trace-Context")
		var traceID, spanID string
		if xcTraceCtx != "" {
			tmp := strings.Split(xcTraceCtx, "/")
			if len(tmp) == 2 {
				traceID = tmp[0]
				spanID = tmp[1]
			}
			ctx = context.WithValue(ctx, log.TraceIDKey{}, traceID)
			ctx = context.WithValue(ctx, log.SpanIDKey{}, spanID)
		}

		lrw := log.NewLogHTTPResponseWriter(w)
		next.ServeHTTP(lrw, r)

		defer func() {
			log.Requestf(ctx, lrw, r)
		}()
	})
}
