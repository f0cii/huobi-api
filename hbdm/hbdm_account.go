package hbdm

import (
	"net/url"
)

// GetAccountInfo 获取用户账户信息
func (c *Client) GetAccountInfo(symbol string) (result AccountInfoResult, err error) {
	path := "/api/v1/contract_account_info"
	params := &url.Values{}
	if symbol != "" {
		params.Add("symbol", symbol)
	}
	_, err = c.doPost(path, params, &result)
	return
}

// GetPositionInfo 用户持仓信息
func (c *Client) GetPositionInfo(symbol string) (result PositionInfoResult, err error) {
	path := "/api/v1/contract_position_info"
	params := &url.Values{}
	if symbol != "" {
		params.Add("symbol", symbol)
	}
	_, err = c.doPost(path, params, &result)
	return
}

// /api/v1/contract_sub_account_list
// /api/v1/contract_sub_account_info
// /api/v1/contract_sub_position_info
// /api/v1/contract_financial_record
// /api/v1/contract_order_limit
// /api/v1/contract_fee
// /api/v1/contract_transfer_limit
// /api/v1/contract_position_limit
// /api/v1/contract_account_position_info
// /api/v1/contract_master_sub_transfer
// /api/v1/contract_master_sub_transfer_record
// /api/v1/contract_api_trading_status
