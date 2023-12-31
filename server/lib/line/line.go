package line

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/emahiro/qrurl/server/lib/jwt"
	"github.com/emahiro/qrurl/server/lib/log"
	"github.com/emahiro/qrurl/server/repository"
)

// singleton
var client *linebot.Client

func NewBot(ctx context.Context, useLongTermToken bool) error {
	var at string
	if !useLongTermToken {
		repo := repository.LineChannelAccessTokenRepository{}
		var err error
		at, err = repo.GetLatestAccessToken(ctx)
		if err != nil {
			return err
		}
		valid, err := CheckIfTokenValid(ctx, at)
		if err != nil {
			return err
		}
		if !valid {
			at, err = PostChannelAccessToken(ctx)
			if err != nil {
				return err
			}
		}
	} else {
		at = os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")
	}

	if err := NewBotClient(at); err != nil {
		return err
	}

	return nil
}

func NewBotClient(at string) error {
	bot, err := linebot.New(os.Getenv("LINE_MESSAGE_CHANNEL_SECRET"), at)
	if err != nil {
		return err
	}

	// set bot to singleton
	client = bot
	return nil
}

func CheckIfTokenValid(ctx context.Context, token string) (bool, error) {
	v := url.Values{}
	v.Add("access_token", token)
	b := strings.NewReader(v.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://api.line.me/oauth2/v2.1/verify", b)
	if err != nil {
		return false, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Infof(ctx, "invalid access-token and regenerate line bot client")
		return false, nil
	}
	return true, nil
}

// PostChannelAccessToken はチャンネルアクセストークンを取得する。
// ChannelAccessToken の登録上限は30件。乱発は禁止。
// ChannelAccessToken の上限に達すると新規の発行はできなくなるので、永続化して都度再利用、有効期限が
// 切れたら refresh する。
// ref: https://developers.line.biz/ja/docs/messaging-api/channel-access-tokens/
func PostChannelAccessToken(ctx context.Context) (string, error) {
	token, err := jwt.CreateToken(ctx)
	if err != nil {
		return "", log.WithStackTracef(err, "failed to create jwt token")
	}
	log.Infof(context.Background(), "token: %s", token)

	form := url.Values{}
	form.Set("grant_type", "client_credentials")
	form.Add("client_assertion_type", "urn:ietf:params:oauth:client-assertion-type:jwt-bearer")
	form.Add("client_assertion", token)
	b := strings.NewReader(form.Encode())

	req, err := http.NewRequest(http.MethodPost, "https://api.line.me/oauth2/v2.1/token", b)
	if err != nil {
		return "", log.WithStackTracef(err, "failed to create request")
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", log.WithStackTracef(err, "failed to request")
	}
	defer resp.Body.Close()

	if resp.StatusCode >= http.StatusBadRequest {
		bb, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", log.WithStackTracef(err, "failed to read response body")
		}
		log.Infof(context.Background(), "token response body. err: %v", string(bb))
		return "", log.WithStackTracef(err, "token response status code is %d body: %v", resp.StatusCode, string(bb))
	}

	var v PostChannelAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return "", log.WithStackTracef(err, "failed to decode response body")
	}

	// persist token process
	// Datastore に取得したアクセストークン、KeyID、有効期限、ClientAssertion を保存する。
	now := time.Now().Unix()
	lineChannelAccessTokenRepo := repository.LineChannelAccessTokenRepository{}
	if err := lineChannelAccessTokenRepo.Create(ctx, repository.LineChannelAccessTokenRepository{
		AccessToken:     v.AccessToken,
		ExpiresIn:       v.ExpiresIn,
		KeyID:           v.KeyID,
		ClientAssertion: token,
		CreatedAt:       now,
		UpdatedAt:       now,
	}); err != nil {
		return "", log.WithStackTracef(err, "failed to persist token")
	}

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
		return nil, log.WithStackTracef(err, "failed to get message content. messageID: %s", messageID)
	}
	defer resp.Content.Close()

	b, err := io.ReadAll(resp.Content)
	if err != nil {
		return nil, log.WithStackTracef(err, "failed to read message content. messageID: %s", messageID)
	}
	return b, nil
}

func ReplyMessage(_ context.Context, replyToken string, text string) error {
	messages := []linebot.SendingMessage{
		linebot.NewTextMessage(text),
	}
	if _, err := client.ReplyMessage(replyToken, messages...).Do(); err != nil {
		return log.WithStackTracef(err, "failed to reply message. replyToken: %s", replyToken)
	}
	return nil
}
