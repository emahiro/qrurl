package middleware

import "net/http"

// Chain is middleware chain for default http handler.
// Inspired by bellow
// - https://zenn.dev/mokmok_dev/articles/30638ae4d15ae6
// - https://github.com/mokmok-dev/middlechain/blob/main/middlechain.go
func Chain(f http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](Chain(f, m[1:cap(m)]...))
}
