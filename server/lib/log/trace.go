package log

import (
	"context"
	"strings"
)

func SetCloudTraceContext(ctx context.Context, xcTraceCtx string) context.Context {
	var traceID, spanID string
	tmp := strings.Split(xcTraceCtx, "/")
	if len(tmp) == 2 {
		traceID = tmp[0]
		spanIDStr := strings.Split(tmp[1], ";")
		if len(spanIDStr) == 2 {
			spanID = spanIDStr[0]
		} else {
			spanID = tmp[1]
		}
	}
	ctx = context.WithValue(ctx, TraceIDKey{}, traceID)
	ctx = context.WithValue(ctx, SpanIDKey{}, spanID)
	return ctx
}
