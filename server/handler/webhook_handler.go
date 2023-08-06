package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/emahiro/qrurl/server/lib"
	"github.com/emahiro/qrurl/server/lib/line"
	"github.com/emahiro/qrurl/server/lib/log"
)

func LineWebHookHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	v := line.LineWebhookRequest{}
	decorder := json.NewDecoder(r.Body)
	for {
		if err := decorder.Decode(&v); err == io.EOF {
			break
		} else if err != nil {
			log.Errorf(ctx, "decode error. err: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	var result string
	for _, event := range v.Events {
		log.Infof(ctx, "event: %v", event)
		message := event.Message
		replyToken := event.ReplyToken

		switch linebot.MessageType(message.Type) {
		case linebot.MessageTypeImage:
			b, err := line.GetMessageContent(ctx, message.Id)
			if err != nil {
				log.Errorf(ctx, "get message content error. err: %v", err)
				if err := line.ReplyMessage(ctx, replyToken, lib.ErrUnknown); err != nil {
					log.Errorf(ctx, "reply message error. err: %v", err)
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusOK)
				return
			}
			log.Infof(ctx, "get message content. %v", b)
			content, err := lib.DecodeQrCode(ctx, b)
			if err != nil {
				log.Errorf(ctx, "decode qr code error. err: %v", err)
				if err := line.ReplyMessage(ctx, replyToken, lib.ErrReadQrCode); err != nil {
					log.Errorf(ctx, "reply message error. err: %v", err)
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
			log.Errorf(ctx, "reply message error. err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	log.Infof(ctx, "result: %v", result)

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(result))
}
