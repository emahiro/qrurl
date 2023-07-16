package middleware

import "net/http"

func Chain(f http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
	// if our chain is done, use the original handlerfunc
	if len(m) == 0 {
		return f
	}
	// otherwise nest the handlerfuncs
	return m[0](Chain(f, m[1:cap(m)]...))
}
