package handler

import (
	"encoding/json"
	"net/http"

	"golang.org/x/exp/slog"

	webhookv1 "github.com/emahiro/qrurl/server/gen/proto/webhook/v1"
)

func LineWebHookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	v := webhookv1.LineWebhookRequest{}
	if err := json.NewDecoder(r.Body).Decode(&v); err != nil {
		slog.ErrorCtx(ctx, "request body parse error", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
