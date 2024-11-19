package open_login

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type QQInfo struct {
	Nickname string `json:"nickname"`
	Gender   string `json:"gender"`
	Avatar   string `json:"figureurl_qq"`
	OpenID   string `json:"open_id"`
}

type QQLogin struct {
	appID     string
	appKey    string
	redirect  string
	code      string
	accessTok string
	openID    string
}

type QQConfig struct {
	AppID    string
	AppKey   string
	Redirect string
}

func NewQQLogin(code string, conf QQConfig) (qqInfo QQInfo, err error) {
	qqLogin := &QQLogin{
		appID:    conf.AppID,
		appKey:   conf.AppKey,
		redirect: conf.Redirect,
		code:     code,
	}
	err = qqLogin.GetAccessToken()
	if err != nil {
		return qqInfo, err
	}
	err = qqLogin.GetOpenID()
	if err != nil {
		return qqInfo, err
	}
	qqInfo, err = qqLogin.GetUserInfo()
	if err != nil {
		return qqInfo, err
	}
	qqInfo.OpenID = qqLogin.openID
	return qqInfo, nil
}

// 获取token
func (q *QQLogin) GetAccessToken() error {
	params := url.Values{}
	params.Add("grant_type", "authorization_code")
	params.Add("client_id", q.appID)
	params.Add("client_secret", q.appKey)
	params.Add("code", q.code)
	params.Add("redirect_uri", q.redirect) // 使用正确的 redirect_uri
	u, err := url.Parse("https://graph.qq.com/oauth2.0/token")
	if err != nil {
		return err
	}
	u.RawQuery = params.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 检查HTTP响应状态
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	qs, err := parseQS(res.Body)
	if err != nil {
		return err
	}

	// 检查是否包含 access_token
	if token, ok := qs["access_token"]; ok && len(token) > 0 {
		q.accessTok = token[0]
	} else {
		return fmt.Errorf("access_token not found")
	}
	return nil
}

func (q *QQLogin) GetOpenID() error {
	// 获取openid
	u, err := url.Parse(fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s", q.accessTok))
	if err != nil {
		return err
	}
	res, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// 检查HTTP响应状态
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	openID, err := getOpenID(res.Body)
	if err != nil {
		return err
	}
	q.openID = openID
	return nil
}

func (q *QQLogin) GetUserInfo() (qqInfo QQInfo, err error) {
	params := url.Values{}
	params.Add("access_token", q.accessTok)
	params.Add("oauth_consumer_key", q.appID)
	params.Add("openid", q.openID) // 使用 q.openID，而不是 q.appID
	u, err := url.Parse("https://graph.qq.com/user/get_user_info")
	if err != nil {
		return qqInfo, err
	}
	u.RawQuery = params.Encode()
	res, err := http.Get(u.String())
	if err != nil {
		return qqInfo, err
	}
	defer res.Body.Close()

	// 检查HTTP响应状态
	if res.StatusCode != http.StatusOK {
		return qqInfo, fmt.Errorf("request failed with status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return qqInfo, err
	}

	err = json.Unmarshal(data, &qqInfo)
	if err != nil {
		return qqInfo, err
	}
	return qqInfo, nil
}

// 将http响应的正文解析为键值对形式
func parseQS(r io.Reader) (val map[string][]string, err error) {
	query, _ := readAll(r)
	val, err = url.ParseQuery(query)
	if err != nil {
		return val, err
	}
	return val, nil
}

// 从http响应的正文中解析出openid
func getOpenID(r io.Reader) (string, error) {
	body, err := readAll(r)
	if err != nil {
		return "", err
	}
	start := strings.Index(body, `"openid":"`) + len(`"openid":"`)
	if start == -1 {
		return "", fmt.Errorf("openid not found")
	}
	end := strings.Index(body[start:], `"`)
	if end == -1 {
		return "", fmt.Errorf("openid not found")
	}
	return body[start : start+end], nil
}

// 读取所有数据并将其转化为字符串
func readAll(r io.Reader) (string, error) {
	buf := new(strings.Builder)
	_, err := io.Copy(buf, r)
	if err != nil {
		return "", fmt.Errorf("error reading from reader: %v", err)
	}
	return buf.String(), nil
}
