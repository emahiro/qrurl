package jwt

import (
	"context"
	"encoding/json"
	"log/slog"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

func CreateToken(ctx context.Context) (string, error) {
	kID := os.Getenv("LINE_PUBLIC_KEY_ID")
	channelID := os.Getenv("LINE_CHANNEL_ID")
	privKey := os.Getenv("LINE_PRIVATE_KEY")

	hdrs := jws.NewHeaders()
	_ = hdrs.Set(jws.AlgorithmKey, "RS256")
	_ = hdrs.Set(jws.TypeKey, "JWT")
	_ = hdrs.Set(jws.KeyIDKey, kID)

	now := time.Now()
	t := jwt.New()
	_ = t.Set(jwt.IssuerKey, channelID)
	_ = t.Set(jwt.SubjectKey, channelID)
	_ = t.Set(jwt.AudienceKey, "https://api.line.me/")
	_ = t.Set(jwt.ExpirationKey, now.Add(30*time.Minute).Unix())
	_ = t.Set("token_exp", 60*60*24*30) // 1month

	buf, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		slog.InfoContext(ctx, "failed to marshal token error", "err", err)
		return "", err
	}

	key, err := jwk.ParseKey([]byte(privKey))
	if err != nil {
		slog.InfoContext(ctx, "failed to parse private key", "err", err)
		return "", err
	}

	signed, err := jws.Sign(buf, jws.WithKey(jwa.RS256, key, jws.WithProtectedHeaders(hdrs)))
	if err != nil {
		slog.InfoContext(ctx, "failed to sign token", "err", err)
		return "", err
	}
	return string(signed[:]), nil
}
