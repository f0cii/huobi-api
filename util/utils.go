package util

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultUserAgent = "Go" // "Mozilla/5.0 (Windows NT 5.1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/31.0.1650.63 Safari/537.36"
)

func HttpGet(client *http.Client, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	return NewHttpRequest(client, "GET", reqUrl, postData, headers)
}

func HttpPost(client *http.Client, reqUrl string, postData string, headers map[string]string) ([]byte, error) {
	return NewHttpRequest(client, "POST", reqUrl, postData, headers)
}

func NewHttpRequest(client *http.Client, method string, reqUrl string, postData string, requestHeaders map[string]string) ([]byte, error) {
	req, _ := http.NewRequest(method, reqUrl, strings.NewReader(postData))
	req.Header.Set("User-Agent", defaultUserAgent)

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
