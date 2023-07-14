package intercepter

import (
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/emahiro/qrurl/server/lib/line"
)

func VerifyChannelAccessToken() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, origReq *http.Request) {
			ctx := origReq.Context()
			// TODO: fetch channel access token from datastore
			var at string
			if at != "" {
				valid, err := line.CheckIfTokenValid(ctx, at)
				if err != nil {
					slog.ErrorCtx(ctx, "failed to check if token is valid: %v", "err=", err)
					http.Error(w, "failed to check if token is valid", http.StatusInternalServerError)
					return
				}
				if !valid {
					if err := line.NewBot(ctx, false); err != nil {
						slog.ErrorCtx(ctx, "failed to fetch new token: %v", "err=", err)
						panic(err)
					}
				}
			}
			next.ServeHTTP(w, origReq)
		})
	}
}
