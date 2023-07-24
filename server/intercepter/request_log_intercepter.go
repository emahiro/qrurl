package intercepter

import (
	"context"

	"github.com/bufbuild/connect-go"

	"github.com/emahiro/qrurl/server/lib/log"
)

func NewRequestLogIntercepter() connect.UnaryInterceptorFunc {
	intercepter := func(next connect.UnaryFunc) connect.UnaryFunc {
		return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			log.Infof(ctx, "this is request info. req: %+v", req)
			log.ConnectRequestf(ctx, req)
			return next(ctx, req)
		}
	}
	return connect.UnaryInterceptorFunc(intercepter)
}
