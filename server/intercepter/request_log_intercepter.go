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
			ctx = context.WithValue(ctx, log.RequestTimeKey{}, time.Now())

			xcTraceCtx := req.Header().Get("X-Cloud-Trace-Context")
			var traceID, spanID string
			if xcTraceCtx != "" {
				tmp := strings.Split(xcTraceCtx, "/")
				if len(tmp) == 2 {
					traceID = tmp[0]
					spanIDStr := strings.Split(tmp[1], ";")
					if len(spanIDStr) == 2 {
						spanID = spanIDStr[0]
					}
				}
				ctx = context.WithValue(ctx, log.TraceIDKey{}, traceID)
				ctx = context.WithValue(ctx, log.SpanIDKey{}, spanID)
			}

			resultStatus := http.StatusOK

			resp, err := next(ctx, req)
			if err != nil {
				if err, ok := err.(*connect.Error); ok {
					switch err.Code() {
					case connect.CodeUnknown:
						resultStatus = http.StatusInternalServerError
					default:
						// error code の分類分け
						resultStatus = http.StatusBadRequest
					}
				}
			}
			// 終了後に Logging する
			defer func() {
				log.ConnectRequestf(ctx, log.ConnectRequestInfo{
					Req:    req,
					Resp:   resp,
					Status: resultStatus,
				})
			}()
			return resp, err
		}
	}
	return connect.UnaryInterceptorFunc(intercepter)
}
