package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"os"
)

var lineMessageChannelSecret = os.Getenv("LINE_MESSAGE_CHANNEL_SECRET")

func VerifyLine(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, origReq *http.Request) {
		ctx := origReq.Context()

		copyReq := origReq.Clone(ctx)
		b, err := io.ReadAll(copyReq.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer copyReq.Body.Close()

		signagure := copyReq.Header.Get("x-line-signature")
		decoded, err := base64.StdEncoding.DecodeString(signagure)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		hash := hmac.New(sha256.New, []byte(lineMessageChannelSecret))
		hash.Write(b)
		verified := hmac.Equal(decoded, hash.Sum(nil))
		if !verified {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		origReq.Body = io.NopCloser(bytes.NewBuffer(b))
	})
}
