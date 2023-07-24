package intercepter

import (
	"context"
	"strings"

	"github.com/bufbuild/connect-go"

	"github.com/emahiro/qrurl/server/lib/log"
)

func NewRequestLogIntercepter() connect.UnaryInterceptorFunc {
	intercepter := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
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
			log.Infof(ctx, "this is request info. req: %+v", req)
			log.ConnectRequestf(ctx, req)
			return next(ctx, req)
		}
	}
	return connect.UnaryInterceptorFunc(intercepter)
}
