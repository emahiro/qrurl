package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/emahiro/qrurl/server/lib"
	"github.com/emahiro/qrurl/server/lib/line"
)

func LineWebHookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	token, err := line.PostChannelAccessToken()
	if err != nil {
		slog.ErrorCtx(ctx, "post channel access token error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	bot, err := line.NewLineBot(token)
	if err != nil {
		slog.ErrorCtx(ctx, "new line bot error", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// MessageID から Messsage の詳細を取得
	// 画像の場合のみ対応
	// 画像を byte に変換
	// lib.DecodeQrCode を呼び出す。
	// URL を取得
	// ユーザーへの応答をする

	v := line.LineWebhookRequest{}
	decorder := json.NewDecoder(r.Body)
	for {
		if err := decorder.Decode(&v); err == io.EOF {
			break
		} else if err != nil {
			slog.ErrorCtx(ctx, "request body parse error", "err", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var result string
	for _, event := range v.Events {
		slog.InfoCtx(ctx, "event", "event", event)
		message := event.Message
		replyToken := event.ReplyToken
		switch message.Type {
		case "text":
			result = message.Text
		case "image":
			b, err := bot.GetMessageContent(ctx, message.Id)
			if err != nil {
				slog.ErrorCtx(ctx, "get message content error", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			slog.InfoCtx(ctx, "get message content", "b", b)
			content, err := lib.DecodeQrCode(ctx, b)
			if err != nil {
				slog.ErrorCtx(ctx, "decode qr code error", "err", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			result = content
		default:
			slog.ErrorCtx(ctx, "not supported type", "type", message.Type)
			result = "not supported"
		}
	}

	slog.InfoCtx(ctx, "result", "result", result)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
