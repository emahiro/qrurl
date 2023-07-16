package repository

import (
	"context"

	"github.com/emahiro/qrurl/server/infra/firestore"
)

type LineChannelAccessTokenRepository struct {
	AccessToken     string `firestore:"access_token,omitempty"`
	TokenType       string `firestore:"token_type,omitempty"`
	ExpiresIn       int64  `firestore:"expires_in,omitempty"`
	KeyID           string `firestore:"key_id,omitempty"`
	ClientAssertion string `firestore:"client_assertion,omitempty"`
	CreatedAt       int64  `firestore:"created_at,omitempty"`
	UpdatedAt       int64  `firestore:"updated_at,omitempty"`
}

const LineChannelAccessTokenCollection = "LineChannelAccessToken"

func (r LineChannelAccessTokenRepository) Create(ctx context.Context, src LineChannelAccessTokenRepository) error {
	return firestore.Add(ctx, LineChannelAccessTokenCollection, src)
}

func (r LineChannelAccessTokenRepository) GetLatestAccessToken(ctx context.Context) (string, error) {
	dst, err := firestore.Query[LineChannelAccessTokenRepository](ctx, LineChannelAccessTokenCollection, firestore.QueryOption{
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
