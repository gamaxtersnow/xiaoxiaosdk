package xiaoxiaosdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/syncx"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const schoolAuthorizationTokenKey = "xiaoxiao:login:token"
const schoolTokenKey = "X-Xiao-Token"
const schoolAuthorizationTokenExpire = 2 * time.Hour

type XiaoxiaoApiConf struct {
	UserName    string
	PassWord    string
	DomainAlias string
	Device      string
	Rate        int   `json:",default=300"`
	Retries     int   `json:",default=3"`
	Delay       int64 `json:",default=500"`
	BaseUrl     string
	CacheConf   cache.CacheConf
}
type HttpClient struct {
	client  *http.Client
	ticker  *time.Ticker
	limiter chan struct{}
	conf    XiaoxiaoApiConf
	cache   cache.Cache
}

type SchoolToken struct {
	BodyToken   string `json:"body_token"`
	HeaderToken string `json:"header_token"`
}

type SchoolTokenResp struct {
	AccountId int64   `json:"accountId"`
	Account   Account `json:"account"`
	TimeStamp int64   `json:"timestamp"`
}

type Account struct {
	Id     int64   `json:"id"`
	Mobile string  `json:"mobile"`
	Money  float64 `json:"money"`
	Token  string  `json:"token"`
}

func NewHttpClient(conf XiaoxiaoApiConf) *HttpClient {
	stats := cache.NewStat("xiaoxiao-sdk")
	errNotFound := errors.New("item not found in xiaoxiao sdk cache")
	cc := cache.New(conf.CacheConf, syncx.NewSingleFlight(), stats, errNotFound)
	ticker := time.NewTicker(time.Minute / time.Duration(conf.Rate))
	limiter := make(chan struct{}, conf.Rate)
	go func() {
		defer close(limiter)
		for range ticker.C {
			limiter <- struct{}{}
		}
	}()
	return &HttpClient{
		client: &http.Client{
			Timeout: 15 * time.Second,
		},
		ticker:  ticker,
		limiter: limiter,
		conf:    conf,
		cache:   cc,
	}
}
func (h *HttpClient) setLoginParams() string {
	params := url.Values{}
	params.Set("username", h.conf.UserName)
	params.Set("password", h.conf.PassWord)
	params.Set("domainalias", h.conf.DomainAlias)
	params.Set("device", h.conf.Device)
	return params.Encode()
}
func (h *HttpClient) setGlobalHeader(request *http.Request) {
	request.Header.Set("content-type", "application/x-www-form-urlencoded")
	request.Header.Set("referer", "https://"+h.conf.DomainAlias+".xiaosaas.com/")
	request.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/129.0.0.0 Safari/537.36")
}
func (h *HttpClient) Token() (*SchoolToken, error) {
	token := &SchoolToken{}
	err := h.cache.Get(schoolAuthorizationTokenKey, token)
	if err == nil {
		return token, nil
	}
	fmt.Println("未获取到token", h.getRequestUrl("/login"))
	request, _ := http.NewRequest(http.MethodPost, h.conf.BaseUrl+"/login", bytes.NewBufferString(h.setLoginParams()))
	h.setGlobalHeader(request)
	resp := &http.Response{}
	resp, err = h.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	headerToken := resp.Header.Get(schoolTokenKey)
	if headerToken == "" {
		return nil, errors.New("header中未获取校校token值")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyToken := &SchoolTokenResp{}
	err = json.Unmarshal(body, bodyToken)
	if err != nil {
		return nil, err
	}
	if bodyToken.Account.Token == "" {
		return token, errors.New("body中未获取校校token值")
	}
	token.HeaderToken = headerToken
	token.BodyToken = bodyToken.Account.Token
	_ = h.cache.SetWithExpire(schoolAuthorizationTokenKey, token, schoolAuthorizationTokenExpire)
	return token, nil
}
func (h *HttpClient) doRequest(method, url string, header map[string]string, body io.Reader) (resp *http.Response, err error) {
	var req *http.Request
	for i := 0; i <= h.conf.Retries; i++ {
		startTime := time.Now()
		<-h.limiter
		req, err = h.NewRequest(method, url, header, body)
		if err != nil {
			break
		}
		resp, err = h.client.Do(req)
		if err == nil {
			break
		}
		if i == h.conf.Retries {
			break
		}
		endTime := time.Now()
		duration := endTime.Sub(startTime)
		milliseconds := duration.Milliseconds()
		if milliseconds < h.conf.Delay {
			sleepDuration := h.conf.Delay - milliseconds
			time.Sleep(time.Duration(sleepDuration) * time.Millisecond)
		}
	}
	return resp, err
}
func (h *HttpClient) Get(path string, params url.Values) (*http.Response, error) {
	resp, err := h.doRequest(http.MethodGet, h.getRequestUrl(path)+"&"+params.Encode(), nil, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *HttpClient) Post(path string, body io.Reader) (*http.Response, error) {
	resp, err := h.doRequest(http.MethodPost, h.getRequestUrl(path), nil, body)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *HttpClient) getRequestUrl(path string) string {
	return strings.TrimRight(h.conf.BaseUrl, "/") + "/" + strings.Trim(path, "/") + "?tok={tok}&lang=cn"
}
func (h *HttpClient) NewRequest(method, requestUrl string, header map[string]string, body io.Reader) (*http.Request, error) {
	xToken, err := h.Token()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, strings.Replace(requestUrl, "{tok}", xToken.BodyToken, 1), body)
	if err != nil {
		return nil, err
	}
	h.setGlobalHeader(req)
	req.Header.Set(schoolTokenKey, xToken.HeaderToken)
	if len(header) > 0 {
		for k, v := range header {
			req.Header.Set(k, v)
		}
	}
	return req, nil
}

func (h *HttpClient) Stop() {
	h.ticker.Stop()
}
