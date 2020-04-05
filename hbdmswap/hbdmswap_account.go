package hbdmswap

import "net/url"

// GetAccountInfo 获取用户账户信息
func (c *Client) GetAccountInfo(contractCode string) (result AccountInfoResult, err error) {
	path := "/swap-api/v1/swap_account_info"
	params := &url.Values{}
	if contractCode != "" {
		params.Add("contract_code", contractCode)
	}
	_, err = c.doPost(path, params, &result)
	return
}

// GetPositionInfo 用户持仓信息
func (c *Client) GetPositionInfo(contractCode string) (result PositionInfoResult, err error) {
	path := "/swap-api/v1/swap_position_info"
	params := &url.Values{}
	if contractCode != "" {
		params.Add("contract_code", contractCode)
	}
	_, err = c.doPost(path, params, &result)
	return
}
