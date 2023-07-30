package middleware

import (
	"net/http"

	"github.com/emahiro/qrurl/server/lib/log"
)

func SetTrace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		xcTraceCtx := r.Header.Get("X-Cloud-Trace-Context")
		if xcTraceCtx != "" {
			ctx = log.SetCloudTraceContext(ctx, xcTraceCtx)
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
