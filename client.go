package chinadep

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/yuhu-tech/chinadep-go-sdk/internal/crypto"
)

type client struct {
	appId     string
	appSecret string

	accessToken  string
	refreshToken string
	expiresAt    time.Time

	// TODO(bootun): add HTTP client
}

func NewClient(appId, appSecret string) *client {
	return &client{
		appId:     appId,
		appSecret: appSecret,
	}
}

func (c *client) ExpiresAt() time.Time {
	return c.expiresAt
}

const (
	// status:"0000", code:1, msg: 成功
	// NOTE: only for token API
	tokenCodeSuccess = 1

	// register API success code
	// status:"ok", code:200, msg: "success"
	registerCodeSuccess = 200
)

type tokenResponse struct {
	Status  string `json:"status"`
	Msg     string `json:"msg"`
	Code    int    `json:"code"`
	Success bool   `json:"success"`
	Data    struct {
		AccessToken  string `json:"accessToken"`
		RefreshToken string `json:"refreshToken"`
		ExpiresIn    int    `json:"expiresIn"`
	} `json:"data"`
}

func (c *client) ApplyToken() error {
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

	if err := c.validateResponseStatusCode(resp.StatusCode); err != nil {
		return fmt.Errorf("failed to validate response status code: %w", err)
	}

	var res tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if res.Code == tokenCodeSuccess {
		c.accessToken = res.Data.AccessToken
		c.refreshToken = res.Data.RefreshToken
		c.expiresAt = time.UnixMilli(time.Now().UnixMilli() + int64(res.Data.ExpiresIn))
		return nil
	} else {
		return fmt.Errorf("failed to apply token: %s", res.Msg)
	}
}

func (c *client) RefreshToken() error {
	timestamp := fmt.Sprintf("%d", time.Now().UnixMilli())
	signOrigin := fmt.Sprintf("appId%sappSecret%srefreshToken%stimestamp%s", c.appId, c.appSecret, c.refreshToken, timestamp)
	sign := fmt.Sprintf("%x", crypto.SumSM3(([]byte(signOrigin))))
	refreshTokenURL := fmt.Sprintf("%s%s", BaseURL, "/platform/api/platform/refreshToken")
	param := url.Values{
		"appId":        []string{c.appId},
		"refreshToken": []string{c.refreshToken},
		"timestamp":    []string{timestamp},
		"sign":         []string{sign},
	}.Encode()
	url := fmt.Sprintf("%s?%s", refreshTokenURL, param)

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

	if err := c.validateResponseStatusCode(resp.StatusCode); err != nil {
		return fmt.Errorf("failed to validate response status code: %w", err)
	}

	var res tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to unmarshal response body: %w", err)
	}

	if res.Code == tokenCodeSuccess {
		c.accessToken = res.Data.AccessToken
		c.refreshToken = res.Data.RefreshToken
		c.expiresAt = time.UnixMilli(time.Now().UnixMilli() + int64(res.Data.ExpiresIn))
		return nil
	} else {
		return fmt.Errorf("failed to refresh token: %s", res.Msg)
	}
}

type MetaverseAssetType int

const (
	// 数字资产
	MetaverseAssetTypeDigitalAsset MetaverseAssetType = 1
	// 数字版权
	MetaverseAssetTypeDigitalRights MetaverseAssetType = 2
	// 其他
	MetaverseAssetTypeOther MetaverseAssetType = 3
)

type AssetsRegisterRequest struct {
	// 必填: 元宇宙资产类型
	// 1=数字资产,2=数字版权,3=其他
	MetaverseAssetType MetaverseAssetType `json:"metaverseAssetType,omitempty"`
	// 必填: 联盟方（链方） id
	ChainId string `json:"chainId,omitempty"`
	// 必填: 合约地址
	ContractAddr string `json:"contractAddr,omitempty"`
	// 系列id
	SeriesId *string `json:"seriesId,omitempty"`
	// 系列名称
	SeriesName *string `json:"seriesName,omitempty"`
	// 系列业务类别(1=数字文创,2=文博衍生,3=品牌营销,4=消费场景,5=产业应用,6=数据知识产权)
	SeriesBizType *int `json:"seriesBizType,omitempty"`
	// 系列封面图地址
	SeriesCoverImgUrl *string `json:"seriesCoverImgUrl,omitempty"`
	// 系列描述
	SeriesDesc *string `json:"seriesDesc,omitempty"`
	// 系列媒体形式(1=图片,2=动图,3=视频,4=音频,5=3D 模型,6=文本,7=其他)
	SeriesMediaType *int `json:"seriesMediaType,omitempty"`
	// 系列创建时间
	CreateTime *int64 `json:"createTime,omitempty"`
	// 资产铸造数量
	MintNumber *int64 `json:"mintNumber,omitempty"`
	// 资产系列元数据链接(metadata)
	MetadataUrl *string `json:"metadataUrl,omitempty"`
	// 发行方身份 id 类型(1=身份证,2=统一社会信用代码)
	IssuerIdType *int `json:"issuerIdType,omitempty"`
	// 发行方身份id
	IssuerId *string `json:"issuerId,omitempty"`
	// 发行方名称
	IssuerName *string `json:"issuerName,omitempty"`
	// 品牌方身份 id 类型(1=身份证,2=统一社会信用代码)
	IpIdType *int `json:"ipIdType,omitempty"`
	// 品牌方身份 id
	IpId *string `json:"ipId,omitempty"`
	// 品牌方名称
	IpName *string `json:"ipName,omitempty"`
	// 流通方式(1=平台挂售,2=领取,3=空投)
	CirculationInfo AssetsRegisterRequestCirculationInfo `json:"circulationInfo,omitempty"`
}

type AssetsRegisterRequestCirculationInfo struct {
	PlatformSelling *int `json:"1,omitempty"`
	Receive         *int `json:"2,omitempty"`
	Airdrop         *int `json:"3,omitempty"`
}

type AssetsRegisterResponse struct {
	Status string                   `json:"status"`
	Code   int                      `json:"code"`
	Msg    string                   `json:"msg"`
	Data   []interface{} `json:"data"`
}

// AssetsRegister 数字资产跨链应用接口。 资产数据以 array 方式批量提交， 每批次最大 100 条
func (c *client) AssetsRegister(body []*AssetsRegisterRequest) error {
	b, err := json.Marshal(body)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	assetsRegisterURL := fmt.Sprintf("%s%s", BaseURL, "/api/v1/assets/register")
	req, err := http.NewRequest(http.MethodPost, assetsRegisterURL, bytes.NewReader(b))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("appid", c.appId)
	req.Header.Set("access-token", c.accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if err := c.validateResponseStatusCode(resp.StatusCode); err != nil {
		return fmt.Errorf("failed to validate response status code: %w", err)
	}
	var res AssetsRegisterResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return fmt.Errorf("failed to decode response body: %w", err)
	}
	if res.Code != registerCodeSuccess {
		return fmt.Errorf("failed to register assets: %s", res.Msg)
	}
	return nil
}

// validateResponseStatusCode validate response status code and return error if necessary
func (c *client) validateResponseStatusCode(code int) error {
	if code < http.StatusContinue {
		return fmt.Errorf("unexpected response status code: %d", code)
	}
	switch code {
	case http.StatusOK:
		return nil
	case http.StatusBadRequest:
		// 提交非法数据
		return fmt.Errorf("bad request: %d", code)
	case http.StatusUnauthorized:
		// 鉴权失败
		return fmt.Errorf("unauthorized: %d", code)
	case http.StatusForbidden:
		// 禁止请求, 请联系管理员添加白名单
		return fmt.Errorf("security policy error: %d", code)
	case http.StatusNotFound:
		// 请求的内容不存在
		return fmt.Errorf("resource not found: %d", code)
	case http.StatusInternalServerError:
		// 服务内部错误
		return fmt.Errorf("internal server error: %d", code)
	default:
		return fmt.Errorf("unknown server error: %d", code)
	}
}
