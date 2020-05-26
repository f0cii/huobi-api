package spot

import (
	"net/url"
	"strconv"
)

// GetMarketDepth 市场深度数据
// depth: 	5，10，20
// _type: 深度的价格聚合度，具体说明见下方 step0，step1，step2，step3，step4，step5
func (c *Client) GetMarketDepth(symbol string, depth int, _type string) (result MarketDepthResult, err error) {
	if _type == "" {
		_type = "step0"
	}
	path := "/market/depth"
	params := &url.Values{}
	params.Add("symbol", symbol)
	params.Add("depth", strconv.Itoa(depth))
	params.Add("type", _type)
	_, err = c.doGet(path, params, &result)
	return
}
