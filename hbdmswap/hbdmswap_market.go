package hbdmswap

import (
	"fmt"
	"net/url"
	"strconv"
)

// GetMarketDepth 获取行情深度数据
// step0-150档 step6-20档
// (150档数据) step0, step1, step2, step3, step4, step5（合并深度1-5）；step0时，不合并深度, (20档数据) step6, step7, step8, step9, step10, step11（合并深度7-11）；step6时，不合并深度
func (c *Client) GetMarketDepth(contractCode string, _type string) (result MarketDepthResult, err error) {
	path := "/swap-ex/market/depth"
	params := &url.Values{}
	params.Add("contract_code", contractCode)
	params.Add("type", _type)
	//var resp []byte
	_, err = c.doGet(path, params, &result)
	//log.Printf("%v", string(resp))
	return
}

// 获取K线数据
// symbol: BTC-USD
// period: 1min, 5min, 15min, 30min, 60min, 4hour, 1day, 1mon

func (c *Client) GetKLine(symbol string, period string, size int, from int64, to int64) (result KLineResult, err error) {
	path := "/swap-ex/market/history/kline"
	params := &url.Values{}
	params.Add("contract_code", symbol)
	params.Add("period", period)
	if size <= 0 {
		size = 150
	}
	params.Add("size", strconv.Itoa(size))
	if from != 0 {
		params.Add("from", fmt.Sprintf("%v", from))
	}
	if to != 0 {
		params.Add("to", fmt.Sprintf("%v", to))
	}
	//var resp []byte
	_, err = c.doGet(path, params, &result)
	//log.Printf("%v", string(resp))
	return
}
