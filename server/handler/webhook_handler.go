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

	// MessageID から Messsage の詳細を取得
	// 画像の場合のみ対応
	// 画像を byte に変換
	// lib.DecodeQrCode を呼び出す。
	// URL を取得
	// ユーザーへの応答をする

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
