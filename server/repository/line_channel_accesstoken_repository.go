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
	dst, err := firestore.GetLatestOne[LineChannelAccessTokenRepository](ctx, LineChannelAccessTokenCollection)
	if err != nil {
		return "", err
	}
	if dst == nil {
		return "", nil
	}
	return dst.AccessToken, nil
}
