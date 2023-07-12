package line

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/v7/linebot"
	"golang.org/x/exp/slog"

	"github.com/emahiro/qrurl/server/lib/jwt"
)

// Described at https://developers.line.biz/ja/reference/messaging-api/#webhook-event-objects
type LineWebhookRequest struct {
	Destination string              `json:"destination"`
	Events      []*LineWebhookEvent `json:"events"`
}

type LineWebhookEvent struct {
	Type    string              `json:"type"`
	Message *LineWebhookMessage `json:"message"`
}

type LineWebhookMessage struct {
	Type           string             `json:"type"`
	Id             string             `json:"id"`
	Text           string             `json:"text"`
	ReplyToken     string             `json:"replyToken"`
	WebhookEventId string             `json:"webhookEventId"`
	Timestamp      int32              `json:"timestamp"`
	Source         *LineWebhookSource `json:"source"`
}

type LineWebhookSource struct {
	Type   string `json:"type"`
	UserId string `json:"userId"`
}

// Described at https://developers.line.biz/ja/reference/messaging-api/#issue-channel-access-token-v2-1
type PostChannelAccessTokenRequest struct {
	GrantType           string `json:"grant_type"`
	ClientAssertionType string `json:"client_assertion_type"`
	ClientAssertion     string `json:"client_assertion"`
}

type PostChannelAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	KeyID       string `json:"key_id"`
}

// PostChannelAccessToken はチャンネルアクセストークンを取得する。
func PostChannelAccessToken() (string, error) {
	token, err := jwt.CreateToken()
	if err != nil {
		return "", err
	}
	slog.InfoCtx(context.Background(), "token", "jwt", token)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", token)
	b := strings.NewReader(form.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://api.line.me/oauth2/v2.1/token", b)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		bb, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", err
		}
		slog.InfoCtx(context.Background(), "token response body", "err", string(bb))
		return "", errors.New(string(bb))
	}

	var v PostChannelAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", err
	}
	return v.AccessToken, nil
}

type LineBot struct {
	bot *linebot.Client
}

func NewLineBot(at string) (*LineBot, error) {
	bot, err := linebot.New(os.Getenv("LINE_MESSAGE_CHANNEL_SECRET"), at)
	if err != nil {
		return nil, err
	}
	return &LineBot{
		bot: bot,
	}, nil
}

func (c *LineBot) GetMessageContent(ctx context.Context, messageID string) ([]byte, error) {
	resp, err := c.bot.GetMessageContent(messageID).Do()
	if err != nil {
		return nil, err
	}
	defer resp.Content.Close()

	b, err := io.ReadAll(resp.Content)
	if err != nil {
		return nil, err
	}
	return b, nil
}
