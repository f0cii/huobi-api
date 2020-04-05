package hbdm

import (
	"fmt"
	"net/url"
)

// Order 合约下单
// https://huobiapi.github.io/docs/dm/v1/cn/#9dc85ffb46
/*
symbol	string	true	"BTC","ETH"...
contract_type	string	true	合约类型 ("this_week":当周 "next_week":下周 "quarter":季度)
contract_code	string	true	BTC180914
client_order_id	long	false	客户自己填写和维护，必须为数字
price	decimal	false	价格
volume	long	true	委托数量(张)
direction	string	true	"buy":买 "sell":卖
offset	string	true	"open":开 "close":平
lever_rate	int	true	杠杆倍数[“开仓”若有10倍多单，就不能再下20倍多单]
order_price_type	string	true	订单报价类型 "limit":限价 "opponent":对手价 "post_only":只做maker单,post only下单只受用户持仓数量限制,optimal_5：最优5档、optimal_10：最优10档、optimal_20：最优20档，ioc:IOC订单，fok：FOK订单, "opponent_ioc"： 对手价-IOC下单，"optimal_5_ioc"：最优5档-IOC下单，"optimal_10_ioc"：最优10档-IOC下单，"optimal_20_ioc"：最优20档-IOC下单,"opponent_fok"： 对手价-FOK下单，"optimal_5_fok"：最优5档-FOK下单，"optimal_10_fok"：最优10档-FOK下单，"optimal_20_fok"：最优20档-FOK下单
*/
func (c *Client) Order(symbol string, contractType string, contractCode string, clientOrderID int64, price float64,
	volume float64, direction string, offset string, leverRate int, orderPriceType string) (result OrderResult, err error) {
	path := "/api/v1/contract_order"
	params := &url.Values{}
	if contractCode != "" {
		params.Add("contract_code", contractCode)
	} else {
		params.Add("symbol", symbol)
		params.Add("contract_type", contractType)
	}
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

// /api/v1/contract_batchorder

/*
 * Cancel 撤销订单
 */
func (c *Client) Cancel(symbol string, orderID int64, clientOrderID int64) (result CancelResult, err error) {
	path := "/api/v1/contract_cancel"
	params := &url.Values{}
	params.Add("symbol", symbol)
	if orderID > 0 {
		params.Add("order_id", fmt.Sprint(orderID))
	}
	if clientOrderID > 0 {
		params.Add("client_order_id", fmt.Sprint(clientOrderID))
	}
	_, err = c.doPost(path, params, &result)
	return
}

// /api/v1/contract_cancelall

/*
 * OrderInfo 获取合约订单信息
 */
func (c *Client) OrderInfo(symbol string, orderID int64, clientOrderID int64) (result OrderInfoResult, err error) {
	path := "/api/v1/contract_order_info"
	params := &url.Values{}
	params.Add("symbol", symbol)
	if orderID > 0 {
		params.Add("order_id", fmt.Sprint(orderID))
	}
	if clientOrderID > 0 {
		params.Add("client_order_id", fmt.Sprint(clientOrderID))
	}
	_, err = c.doPost(path, params, &result)
	return
}

// TODO: OrderDetail
func (c *Client) OrderDetail(symbol string, orderID int64, createdAt int64, orderType int, pageIndex int, pageSize int) (err error) {
	path := "/api/v1/contract_order_detail"
	params := &url.Values{}
	params.Add("symbol", symbol)
	params.Add("order_id", fmt.Sprint(orderID))
	if createdAt > 0 {
		params.Add("created_at", fmt.Sprint(createdAt))
	}
	if orderType > 0 {
		params.Add("order_type", fmt.Sprint(orderType))
	}
	if pageIndex > 0 {
		params.Add("page_index", fmt.Sprint(pageIndex))
	}
	if pageSize > 0 {
		params.Add("page_size", fmt.Sprint(pageSize))
	}
	_, err = c.doPost(path, params, nil)
	return
}

// GetOpenOrders 获取合约当前未成交委托单
// page_index: 1,2,3...
func (c *Client) GetOpenOrders(symbol string, pageIndex int, pageSize int) (result OpenOrdersResult, err error) {
	path := "/api/v1/contract_openorders"
	params := &url.Values{}
	params.Add("symbol", symbol)
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
func (c *Client) GetHisOrders(symbol string, tradeType int, _type int, status int, createDate int,
	pageIndex int, pageSize int, contractCode string, orderType string) (result HisOrdersResult, err error) {
	path := "/api/v1/contract_hisorders"
	params := &url.Values{}
	params.Add("symbol", symbol)
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
	if contractCode != "" {
		params.Add("contract_code", contractCode)
	}
	if orderType != "" {
		params.Add("order_type", orderType)
	}
	_, err = c.doPost(path, params, &result)
	return
}

// /api/v1/contract_matchresults

// LightningClosePosition 闪电平仓下单
func (c *Client) LightningClosePosition(symbol string, contractType string, contractCode string, volume int,
	direction string, clientOrderID int64, orderPriceType string) (result LightningClosePositionResult, err error) {
	path := "/api/v1/lightning_close_position"
	params := &url.Values{}
	if contractCode != "" {
		params.Add("contract_code", contractCode)
	} else {
		params.Add("symbol", symbol)
		params.Add("contract_type", contractType)
	}
	params.Add("volume", fmt.Sprint(volume))
	params.Add("direction", direction)
	if clientOrderID > 0 {
		params.Add("client_order_id", fmt.Sprint(clientOrderID))
	}
	if orderPriceType != "" {
		params.Add("order_price_type", orderPriceType)
	}
	_, err = c.doPost(path, params, &result)
	return
}

// /api/v1/contract_trigger_order
// /api/v1/contract_trigger_cancel
// /api/v1/contract_trigger_cancelall
// /api/v1/contract_trigger_openorders
// /api/v1/contract_trigger_hisorders
