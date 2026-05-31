package repository

import (
	"context"

	"github.com/emahiro/qrurl/server/infra/firestore"
	"github.com/emahiro/qrurl/server/model"
)

type LineChannelAccessTokenRepository struct{}

const LineChannelAccessTokenCollection = "LineChannelAccessToken"

func (r LineChannelAccessTokenRepository) Create(ctx context.Context, src model.LineChannelAccessToken) error {
	return firestore.Add(ctx, LineChannelAccessTokenCollection, src)
}

func (r LineChannelAccessTokenRepository) GetLatestAccessToken(ctx context.Context) (string, error) {
	dst, err := firestore.Query[model.LineChannelAccessToken](ctx, LineChannelAccessTokenCollection, firestore.QueryOption{
		OrderBy: "created_at",
		Desc:    true,
		Limit:   1,
	})
	if err != nil {
		return "", err
	}
	if len(dst) == 0 {
		return "", nil
	}
	return dst[0].AccessToken, nil
}
