package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"golang.org/x/exp/slog"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/emahiro/qrurl/server/lib"
	"github.com/emahiro/qrurl/server/lib/line"
)

func LineWebHookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

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

		switch linebot.MessageType(message.Type) {
		case linebot.MessageTypeText:
			result = message.Text
		case linebot.MessageTypeImage:
			b, err := line.GetMessageContent(ctx, message.Id)
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

		// reply message
		if err := line.ReplyMessage(ctx, replyToken, result); err != nil {
			slog.ErrorCtx(ctx, "reply message error", "err", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	slog.InfoCtx(ctx, "result", "result", result)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
