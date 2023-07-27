package middleware

import "net/http"

func RequestLog(http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, origReq *http.Request) {

	})
}
