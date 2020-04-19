package hbdm

import (
	"fmt"
	"github.com/frankrap/huobi-api/util"
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
	HttpClient         *http.Client
	ProxyURL           string
}

type Client struct {
	params     *ApiParameter
	domain     string
	httpClient *http.Client
}

func (c *Client) Heartbeat() (result HeartbeatResult, err error) {
	var resp []byte
	resp, err = util.HttpGet(c.httpClient, "https://www.hbdm.com/heartbeat/", "", nil)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &result)
	return
}

func (c *Client) doGet(path string, params *url.Values, result interface{}) (resp []byte, err error) {
	url := c.params.Url + path + "?" + params.Encode()
	resp, err = util.HttpGet(
		c.httpClient,
		url,
		"",
		nil,
	)
	if err != nil {
		return
	}

	if c.params.Debug {
		log.Println(string(resp))
	}

	if result == nil {
		return
	}

	err = json.Unmarshal(resp, result)
	return
}

func (c *Client) doPost(path string, params *url.Values, result interface{}) (resp []byte, err error) {
	c.sign("POST", path, params)
	jsonD, _ := util.ValuesToJson(*params)

	url := c.params.Url + path + "?" + params.Encode()
	resp, err = util.HttpPost(
		c.httpClient,
		url,
		string(jsonD),
		defaultPostHeaders,
	)
	if err != nil {
		return
	}

	if c.params.Debug {
		log.Println(string(resp))
	}

	if result == nil {
		return
	}

	err = json.Unmarshal(resp, result)
	return
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
		postForm.Encode())
	signature, _ := util.GetParamHmacSHA256Base64Sign(c.params.SecretKey, payload)
	postForm.Set("Signature", signature)
	return nil
}

func NewClient(params *ApiParameter) *Client {
	domain := strings.Replace(params.Url, "https://", "", -1)
	httpClient := params.HttpClient
	if httpClient == nil {
		transport := util.CloneDefaultTransport()
		if params.ProxyURL != "" {
			transport.Proxy, _ = util.ParseProxy(params.ProxyURL)
		}
		httpClient = &http.Client{
			Timeout:   10 * time.Second,
			Transport: transport,
		}
	}
	return &Client{
		params:     params,
		domain:     domain,
		httpClient: httpClient,
	}
}
