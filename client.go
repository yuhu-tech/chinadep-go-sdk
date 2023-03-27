package chinadep

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/yuhu-tech/chinadep-go-sdk/internal/crypto"
)

type Client struct {
	appId     string
	appSecret string

	accessToken  string
	refreshToken string
	expiresIn    int
}

func NewClient(appId, appSecret string) *Client {
	return &Client{
		appId:     appId,
		appSecret: appSecret,
	}
}

type applyTokenResponse struct {
	Status string `json:"status"`
	Msg    string `json:"msg"`
	Data   struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn    int    `json:"expiresIn"`
	} `json:"data"`
	Code    int  `json:"code"`
	Success bool `json:"success"`
}

func (c *Client) ApplyToken() error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	appId := c.appId
	appSecret := c.appSecret
	signOrigin := fmt.Sprintf("appId%sappSecret%vtimestamp%v", appId, appSecret, timestamp)
	sign := fmt.Sprintf("%x", crypto.SumSM3(([]byte(signOrigin))))
	applyTokenURL := `https://chainbridge.chinadep.com/platform/api/platform/token`
	param := url.Values{
		"appId":     []string{appId},
		"timestamp": []string{timestamp},
		"sign":      []string{sign},
	}.Encode()
	url := fmt.Sprintf("%s?%s", applyTokenURL, param)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		break
	case http.StatusForbidden:
		// Security policy error, please contact the chinadep administrator
		return fmt.Errorf("security policy error: %s", resp.Status)
	default:
		return fmt.Errorf("failed to apply token: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	var res applyTokenResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	// res.Status
	// 0000 成功
	// 1001 错误的appId
	// 1002 错误的sign
	if res.Status == "0000" && res.Code == 1 && res.Msg == "success" && res.Success {
		c.accessToken = res.Data.AccessToken
		c.refreshToken = res.Data.RefreshToken
		c.expiresIn = res.Data.ExpiresIn
		return nil
	} else {
		return fmt.Errorf("failed to apply token: %s", res.Msg)
	}
}
