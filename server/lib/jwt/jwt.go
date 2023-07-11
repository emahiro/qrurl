package jwt

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/lestrrat-go/jwx/v2/jwa"
	"github.com/lestrrat-go/jwx/v2/jwk"
	"github.com/lestrrat-go/jwx/v2/jws"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

var (
	kID        = os.Getenv("LINE_PUBLIC_KEY_ID")
	channelID  = os.Getenv("LINE_CHANNEL_ID")
	pribateKey = os.Getenv("LINE_PRIVATE_KEY")
)

func CreateToken() (string, error) {
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
	_ = t.Set("token_exp", now.Add(24*time.Hour).Unix())

	buf, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		fmt.Printf("failed to marshal token: %s\n", err)
		return "", err
	}

	key, err := jwk.ParseKey([]byte(pribateKey))
	if err != nil {
		fmt.Printf("failed to parse private key: %s\n", err)
		return "", err
	}

	signed, err := jws.Sign(buf, jws.WithKey(jwa.RS256, key, jws.WithProtectedHeaders(hdrs)))
	if err != nil {
		fmt.Printf("failed to sign token: %s\n", err)
		return "", err
	}
	return string(signed[:]), nil
}