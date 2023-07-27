package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/emahiro/qrurl/server/lib/log"
)

func RequestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		requestTime := time.Now()
		ctx = context.WithValue(ctx, log.RequestTimeKey{}, requestTime)

		lrw := log.NewLogHTTPResponseWriter(w)
		next.ServeHTTP(lrw, r)

		defer func() {
			log.Requestf(ctx, lrw, r)
		}()
	})
}
