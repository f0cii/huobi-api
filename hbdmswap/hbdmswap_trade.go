package hbdmswap

import (
	"fmt"
	"net/url"
)

// Order 合约下单
/*
contract_code: BTC-USD
direction: "buy":买 "sell":卖
offset	string	true	"open":开 "close":平
orderPriceType: 订单报价类型 "limit":限价 "opponent":对手价 "post_only":只做maker单,post only下单只受用户持仓数量限制,optimal_5：最优5档、optimal_10：最优10档、optimal_20：最优20档，"fok":FOK订单，"ioc":IOC订单, opponent_ioc"： 对手价-IOC下单，"optimal_5_ioc"：最优5档-IOC下单，"optimal_10_ioc"：最优10档-IOC下单，"optimal_20_ioc"：最优20档-IOC下单,"opponent_fok"： 对手价-FOK下单，"optimal_5_fok"：最优5档-FOK下单，"optimal_10_fok"：最优10档-FOK下单，"optimal_20_fok"：最优20档-FOK下单
*/
func (c *Client) Order(contractCode string, clientOrderID int64, price float64,
	volume float64, direction string, offset string, leverRate int, orderPriceType string) (result OrderResult, err error) {
	path := "/swap-api/v1/swap_order"
	params := &url.Values{}
	params.Add("contract_code", contractCode)
	if clientOrderID > 0 {
		params.Add("client_order_id", fmt.Sprint(clientOrderID))
	}
	if price > 0 {
		params.Add("price", fmt.Sprint(price))
	}
	params.Add("volume", fmt.Sprint(volume))
	params.Add("direction", direction)
	params.Add("offset", offset)
	params.Add("lever_rate", fmt.Sprint(leverRate))
	params.Add("order_price_type", orderPriceType)
	_, err = c.doPost(path, params, &result)
	return
}

/*
 * Cancel 撤销订单
 */
func (c *Client) Cancel(contractCode string, orderID int64, clientOrderID int64) (result CancelResult, err error) {
	path := "/swap-api/v1/swap_cancel"
	params := &url.Values{}
	params.Add("contract_code", contractCode)
	if orderID > 0 {
		params.Add("order_id", fmt.Sprint(orderID))
	}
	if clientOrderID > 0 {
		params.Add("client_order_id", fmt.Sprint(clientOrderID))
	}
	_, err = c.doPost(path, params, &result)
	return
}

/*
 * OrderInfo 获取合约订单信息
 */
func (c *Client) OrderInfo(contractCode string, orderID int64, clientOrderID int64) (result OrderInfoResult, err error) {
	path := "/swap-api/v1/swap_order_info"
	params := &url.Values{}
	params.Add("contract_code", contractCode)
	if orderID > 0 {
		params.Add("order_id", fmt.Sprint(orderID))
	}
	if clientOrderID > 0 {
		params.Add("client_order_id", fmt.Sprint(clientOrderID))
	}
	_, err = c.doPost(path, params, &result)
	return
}

// GetOpenOrders 获取合约当前未成交委托单
// page_index: 1,2,3...
func (c *Client) GetOpenOrders(contractCode string, pageIndex int, pageSize int) (result OpenOrdersResult, err error) {
	path := "/swap-api/v1/swap_openorders"
	params := &url.Values{}
	params.Add("contract_code", contractCode)
	if pageIndex > 0 {
		params.Add("page_index", fmt.Sprint(pageIndex))
	}
	if pageSize > 0 {
		params.Add("page_size", fmt.Sprint(pageSize))
	}
	_, err = c.doPost(path, params, &result)
	return
}

// GetHisOrders 获取合约历史委托
func (c *Client) GetHisOrders(contractCode string, tradeType int, _type int, status int, createDate int,
	pageIndex int, pageSize int) (result HisOrdersResult, err error) {
	path := "/swap-api/v1/swap_hisorders"
	params := &url.Values{}
	params.Add("contract_code", contractCode)
	params.Add("trade_type", fmt.Sprint(tradeType))
	params.Add("type", fmt.Sprint(_type))
	params.Add("status", fmt.Sprint(status))
	params.Add("create_date", fmt.Sprint(createDate))
	if pageIndex > 0 {
		params.Add("page_index", fmt.Sprint(pageIndex))
	}
	if pageSize > 0 {
		params.Add("page_size", fmt.Sprint(pageSize))
	}
	_, err = c.doPost(path, params, &result)
	return
}
