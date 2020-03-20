package hbdm

import (
	"fmt"
	"net/url"
	"strconv"
)

/**
 * 期货行情 https://huobiapi.github.io/docs/dm/v1/cn/#0737c93bf7
 *
 * @param symbol
 *            "BTC","ETH"...
 * @param contractType
 *            合约类型: this_week:当周 next_week:下周 quarter:季度
 * @param contract_code
 *            合约code: BTC200320
 * @return
 */
func (c *Client) GetContractInfo(symbol, contractType, contractCode string) (result ContractInfoResult, err error) {
	path := "/api/v1/contract_contract_info"
	params := &url.Values{}
	if symbol != "" {
		params.Add("symbol", symbol)
	}
	if contractType != "" {
		params.Add("contract_type", contractType)
	}
	if contractCode != "" {
		params.Add("contractCode", contractCode)
	}
	err = c.doGet(path, params, &result)
	return
}

/**
 * 获取合约指数信息 https://huobiapi.github.io/docs/dm/v1/cn/#1028ab8392
 */
func (c *Client) GetContractIndex(symbol string) (result ContractIndexResult, err error) {
	path := "/api/v1/contract_index"
	params := &url.Values{}
	params.Add("symbol", symbol)
	err = c.doGet(path, params, &result)
	return
}

// api/v1/contract_price_limit
// api/v1/contract_open_interest
// api/v1/contract_delivery_price
// api/v1/contract_api_state

// GetMarketDepth 获取行情深度数据
// 如"BTC_CW"表示BTC当周合约，"BTC_NW"表示BTC次周合约，"BTC_CQ"表示BTC季度合约
// 获得150档深度数据，使用step0, step1, step2, step3, step4, step5（step1至step5是进行了深度合并后的深度），使用step0时，不合并深度获取150档数据;获得20档深度数据，使用 step6, step7, step8, step9, step10, step11（step7至step11是进行了深度合并后的深度），使用step6时，不合并深度获取20档数据
// https://api.hbdm.com/market/depth?symbol=BTC_CQ&type=step5
func (c *Client) GetMarketDepth(symbol string, _type string) (result MarketDepthResult, err error) {
	path := "/market/depth"
	params := &url.Values{}
	params.Add("symbol", symbol)
	params.Add("type", _type)
	err = c.doGet(path, params, &result)
	return
}

// 获取K线数据 https://huobiapi.github.io/docs/dm/v1/cn/#k
// symbol	true	string	合约名称		如"BTC_CW"表示BTC当周合约，"BTC_NW"表示BTC次周合约，"BTC_CQ"表示BTC季度合约
// period	true	string	K线类型		1min, 5min, 15min, 30min, 60min,4hour,1day, 1mon
// size	false	integer	获取数量	150	[1,2000]
// from	false	integer	开始时间戳 10位 单位S
// to	false	integer	结束时间戳 10位 单位S
func (c *Client) GetKLine(symbol string, period string, size int, from int64, to int64) (result KLineResult, err error) {
	path := "/market/history/kline"
	params := &url.Values{}
	params.Add("symbol", symbol)
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
	err = c.doGet(path, params, &result)
	//log.Printf("%v", string(resp))
	return
}

// /market/detail/merged
// /market/trade
// /market/history/trade
// /api/v1/contract_risk_info
// /api/v1/contract_insurance_fund
// /api/v1/contract_adjustfactor
// /api/v1/contract_his_open_interest
// /api/v1/contract_elite_account_ratio
// /api/v1/contract_elite_position_ratio
// /api/v1/contract_liquidation_orders
// /api/v1/index/market/history/index
// /api/v1/index/market/history/basis
