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

// singleton
var client *linebot.Client

func NewBot(ctx context.Context, useLongTermToken bool) error {
	// [TODO]: checking if latest accesstoken is valid
	// 1. fetch token from datastore.
	// 2. check if token is valid.
	// 3. if valid, use it.
	// 4. if not valid, fetch new token from LINE API or using long term token.

	at := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	if !useLongTermToken || at != "" {
		// check validation
		t, err := postChannelAccessToken(ctx)
		if err != nil {
			return err
		}
		at = t
	}

	bot, err := linebot.New(os.Getenv("LINE_MESSAGE_CHANNEL_SECRET"), at)
	if err != nil {
		return err
	}

	client = bot
	return nil
}

// PostChannelAccessToken はチャンネルアクセストークンを取得する。
// ChannelAccessToken の登録上限は30件。乱発は禁止。
// ChannelAccessToken の上限に達すると新規の発行はできなくなるので、永続化して都度再利用、有効期限が
// 切れたら refresh する。
// ref: https://developers.line.biz/ja/docs/messaging-api/channel-access-tokens/
func postChannelAccessToken(ctx context.Context) (string, error) {
	token, err := jwt.CreateToken(ctx)
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

	// persist token process
	// Datastore に取得したアクセストークン、KeyID、有効期限、ClientAssertion を保存する。

	return v.AccessToken, nil
}

// Described at https://developers.line.biz/ja/reference/messaging-api/#webhook-event-objects
type LineWebhookRequest struct {
	Destination string              `json:"destination"`
	Events      []*LineWebhookEvent `json:"events"`
}

type LineWebhookEvent struct {
	Type           string             `json:"type"`
	Message        *LineMessage       `json:"message"`
	ReplyToken     string             `json:"replyToken"`
	WebhookEventId string             `json:"webhookEventId"`
	Timestamp      int64              `json:"timestamp"`
	Source         *LineWebhookSource `json:"source"`
}

type LineMessage struct {
	Type string `json:"type"`
	Id   string `json:"id"`
	Text string `json:"text"`
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

func GetMessageContent(_ context.Context, messageID string) ([]byte, error) {
	resp, err := client.GetMessageContent(messageID).Do()
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

func ReplyMessage(_ context.Context, replyToken string, text string) error {
	messages := []linebot.SendingMessage{
		linebot.NewTextMessage(text),
	}
	if _, err := client.ReplyMessage(replyToken, messages...).Do(); err != nil {
		return err
	}
	return nil
}
