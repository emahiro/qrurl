package intercepter

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/bufbuild/connect-go"

	"github.com/emahiro/qrurl/server/lib/log"
)

func NewRequestLogIntercepter() connect.UnaryInterceptorFunc {
	intercepter := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			requestTime := time.Now()

			xcTraceCtx := req.Header().Get("X-Cloud-Trace-Context")
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

			resp, err := next(ctx, req)
			if err != nil {
				log.ConnectRequestf(ctx, log.ConnectRequestInfo{
					Req:         req,
					Resp:        nil,
					Status:      http.StatusInternalServerError,
					RequestTime: requestTime,
					Duration:    time.Since(requestTime),
				})
				return nil, err
			}
			defer func() {
				log.ConnectRequestf(ctx, log.ConnectRequestInfo{
					Req:         req,
					Resp:        resp,
					Status:      http.StatusOK,
					RequestTime: requestTime,
					Duration:    time.Since(requestTime),
				})
			}()
			return resp, nil
		}
	}
	return connect.UnaryInterceptorFunc(intercepter)
}
