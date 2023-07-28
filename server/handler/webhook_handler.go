package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/exp/slog"

	"github.com/emahiro/qrurl/server/lib"
	"github.com/emahiro/qrurl/server/lib/line"
)

func LineWebHookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fmt.Println("this is line webhook handler")

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
		case linebot.MessageTypeImage:
			b, err := line.GetMessageContent(ctx, message.Id)
			if err != nil {
				slog.ErrorCtx(ctx, "get message content error", "err", err)
				if err := line.ReplyMessage(ctx, replyToken, lib.ErrUnknown); err != nil {
					slog.ErrorCtx(ctx, "reply message error", "err", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
			slog.InfoCtx(ctx, "get message content", "b", b)
			content, err := lib.DecodeQrCode(ctx, b)
			if err != nil {
				slog.ErrorCtx(ctx, "decode qr code error", "err", err)
				if err := line.ReplyMessage(ctx, replyToken, lib.ErrReadQrCode); err != nil {
					slog.ErrorCtx(ctx, "reply message error", "err", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
			result = content
		default:
			result = lib.ErrNotSupportedMediaType
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
