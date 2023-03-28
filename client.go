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
	expiresAt    time.Time
}

func NewClient(appId, appSecret string) *Client {
	return &Client{
		appId:     appId,
		appSecret: appSecret,
	}
}

func (c *Client) ExpiresAt() time.Time {
	return c.expiresAt
}

const (
	// code:0000, msg: 成功
	statusSuccess = "0000"
	// code:1001, msg: 错误的appId
	statusBadAppid = "1001"
	// code:1002, msg: 错误的sign
	statusBadSign = "1002"
	// code:1003, msg: 错误的refreshToken
	statusBadRefreshToken = "1003"
	// code:1004, msg: 错误的accessToken
	statusBadAccessToken = "1004"
	// code:1005, msg: refreshToken超时
	statusRefreshTokenTimeout = "1005"
	// code:1006, msg: accessToken超时
	statusAccessTokenTimeout = "1006"
)

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
	signOrigin := fmt.Sprintf("appId%sappSecret%vtimestamp%v", c.appId, c.appSecret, timestamp)
	sign := fmt.Sprintf("%x", crypto.SumSM3(([]byte(signOrigin))))
	applyTokenURL := fmt.Sprintf("%s%s", BaseURL, "/platform/api/platform/token")
	param := url.Values{
		"appId":     []string{c.appId},
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

	if res.Status == statusSuccess {
		c.accessToken = res.Data.AccessToken
		c.refreshToken = res.Data.RefreshToken
		c.expiresAt = time.UnixMilli(time.Now().UnixMilli() + int64(res.Data.ExpiresIn))
		return nil
	} else {
		return fmt.Errorf("failed to apply token: %s", res.Msg)
	}
}

func (c *Client) RefreshToken() error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signOrigin := fmt.Sprintf("appId%sappSecret%srefreshToken%stimestamp%s", c.appId, c.appSecret, c.refreshToken, timestamp)
	sign := fmt.Sprintf("%x", crypto.SumSM3(([]byte(signOrigin))))
	applyTokenURL := fmt.Sprintf("%s%s",BaseURL, "/platform/api/platform/refreshToken")
	param := url.Values{
		"appId":        []string{c.appId},
		"refreshToken": []string{c.refreshToken},
		"timestamp":    []string{timestamp},
		"sign":         []string{sign},
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
		return fmt.Errorf("failed to refresh token: %s", resp.Status)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}
	var res applyTokenResponse
	if err := json.Unmarshal(b, &res); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if res.Status == statusSuccess {
		c.accessToken = res.Data.AccessToken
		c.refreshToken = res.Data.RefreshToken
		c.expiresAt = time.UnixMilli(time.Now().UnixMilli() + int64(res.Data.ExpiresIn))
		return nil
	} else {
		return fmt.Errorf("failed to refresh token: %s", res.Msg)
	}
}
