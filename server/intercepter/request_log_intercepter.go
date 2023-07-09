package intercepter

import (
	"context"
	"fmt"

	"github.com/bufbuild/connect-go"
	"golang.org/x/exp/slog"
)

func NewRequestLogIntercepter() connect.UnaryInterceptorFunc {
	intercepter := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			slog.InfoCtx(ctx, "request", fmt.Sprintf("%+v", req))
			return next(ctx, req)
		}
	}
	return connect.UnaryInterceptorFunc(intercepter)
}
