package utils

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	DefaultTransport *http.Transport
)

func init() {
	DefaultTransport = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   5 * time.Second,
			KeepAlive: 30 * time.Second,
			DualStack: true,
		}).DialContext,
		IdleConnTimeout:        1 * time.Minute,
		TLSHandshakeTimeout:    10 * time.Second,
		ExpectContinueTimeout:  1 * time.Second,
		DisableKeepAlives:      false,
		MaxResponseHeaderBytes: 1 << 15,
	}
}

func DefaultHttpClient(proxyURL string) *http.Client {
	transport := CloneDefaultTransport()
	if proxyURL != "" {
		transport.Proxy, _ = ParseProxy(proxyURL)
	}
	httpClient := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	return httpClient
}

func HttpGet(client *http.Client, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	return NewHttpRequest(client, "GET", reqUrl, postData, headers)
}

func HttpPost(client *http.Client, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	return NewHttpRequest(client, "POST", reqUrl, postData, headers)
}

func NewHttpRequest(client *http.Client, method string, reqUrl string, postData string, requestHeaders map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(method, reqUrl, strings.NewReader(postData))
	//req.Header.Set("User-Agent", defaultUserAgent)

	if requestHeaders != nil {
		for k, v := range requestHeaders {
			req.Header.Add(k, v)
		}
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	bodyData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("HttpStatusCode: %d, Desc: %s", resp.StatusCode, string(bodyData)))
	}

	return bodyData, nil
}

func ValuesToJson(v url.Values) ([]byte, error) {
	m := make(map[string]interface{})
	for k, vv := range v {
		if len(vv) == 1 {
			m[k] = vv[0]
		} else {
			m[k] = vv
		}
	}
	return json.Marshal(m)
}

func GetParamHmacSHA256Base64Sign(secret, params string) (string, error) {
	mac := hmac.New(sha256.New, []byte(secret))
	_, err := mac.Write([]byte(params))
	if err != nil {
		return "", err
	}
	signByte := mac.Sum(nil)
	return base64.StdEncoding.EncodeToString(signByte), nil
}

func CloneDefaultTransport() *http.Transport {
	return &http.Transport{DialContext: DefaultTransport.DialContext,
		IdleConnTimeout:        DefaultTransport.IdleConnTimeout,
		TLSHandshakeTimeout:    DefaultTransport.TLSHandshakeTimeout,
		ExpectContinueTimeout:  DefaultTransport.ExpectContinueTimeout,
		MaxResponseHeaderBytes: DefaultTransport.MaxResponseHeaderBytes,
		Proxy:                  DefaultTransport.Proxy,
		DisableKeepAlives:      DefaultTransport.DisableKeepAlives,
	}
}

// "socks5://127.0.0.1:1080"
func ParseProxy(proxyURL string) (res func(*http.Request) (*url.URL, error), err error) {
	var purl *url.URL
	purl, err = url.Parse(proxyURL)
	if err != nil {
		return
	}
	res = http.ProxyURL(purl)
	return
}
