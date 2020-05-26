package spot

import (
	"github.com/frankrap/huobi-api/utils"
	"log"
	"net/http"
	"net/url"
)

type ApiParameter struct {
	Debug              bool
	AccessKey          string
	SecretKey          string
	EnablePrivateSign  bool
	BaseURL            string
	PrivateKeyPrime256 string
	HttpClient         *http.Client
	ProxyURL           string
}

type Client struct {
	params     *ApiParameter
	httpClient *http.Client
}

func (c *Client) doGet(path string, params *url.Values, result interface{}) (resp []byte, err error) {
	url := c.params.BaseURL + path + "?" + params.Encode()
	resp, err = utils.HttpGet(
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

func NewClient(params *ApiParameter) *Client {
	httpClient := params.HttpClient
	if httpClient == nil {
		httpClient = utils.DefaultHttpClient("")
	}
	return &Client{
		httpClient: httpClient,
		params:     params,
	}
}
