package hbdm

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var (
	defaultPostHeaders = map[string]string{"Content-Type": "application/json", "Accept-Language": "zh-cn"}
)

type ApiParameter struct {
	Debug              bool
	AccessKey          string
	SecretKey          string
	EnablePrivateSign  bool
	Url                string
	PrivateKeyPrime256 string
}

type Client struct {
	params     *ApiParameter
	domain     string
	httpClient *http.Client
}

func (c *Client) doGet(path string, params *url.Values, result interface{}) (err error) {
	url := c.params.Url + path + "?" + params.Encode()
	var resp []byte
	resp, err = HttpGet(c.httpClient, url, "", nil)
	if err != nil {
		return
	}
	if result == nil {
		return
	}
	err = json.Unmarshal(resp, result)
	return
}

func (c *Client) doPost(path string, params *url.Values, result interface{}) error {
	c.sign("POST", path, params)
	jsonD, _ := ValuesToJson(*params)

	url := c.params.Url + path + "?" + params.Encode()
	resp, err := HttpPost(c.httpClient,
		url,
		string(jsonD),
		defaultPostHeaders,
	)

	if err != nil {
		return err
	}

	if c.params.Debug {
		log.Println(string(resp))
	}

	if result == nil {
		return nil
	}

	err = json.Unmarshal(resp, result)
	return err
}

func (c *Client) sign(reqMethod, path string, postForm *url.Values) error {
	postForm.Set("AccessKeyId", c.params.AccessKey)
	postForm.Set("SignatureMethod", "HmacSHA256")
	postForm.Set("SignatureVersion", "2")
	postForm.Set("Timestamp", time.Now().UTC().Format("2006-01-02T15:04:05"))
	payload := fmt.Sprintf("%s\n%s\n%s\n%s",
		reqMethod,
		c.domain,
		path,
		postForm.Encode(),
	)
	signature, _ := GetParamHmacSHA256Base64Sign(c.params.SecretKey, payload)
	postForm.Set("Signature", signature)
	return nil
}

func NewClient(params *ApiParameter) *Client {
	domain := strings.Replace(params.Url, "https://", "", -1)
	return &Client{
		params:     params,
		domain:     domain,
		httpClient: &http.Client{},
	}
}
