# huobi-api
Huobi api for Golang. 火币交割合约和永续合约接口.

[![](https://img.shields.io/badge/api-huobi-blue.svg)](https://huobiapi.github.io/docs/dm/v1/cn/)

An implementation of [Huobi-DM API](https://huobiapi.github.io/docs/dm/v1/cn/) and [Huobi-DM-Swap API](https://docs.huobigroup.com/docs/coin_margined_swap/v1/cn/).

## Installation
```
go get github.com/frankrap/huobi-api
```

## Usage
```go
package main

import (
	"github.com/frankrap/huobi-api/hbdm"
	"log"
)

func main() {
	accessKey := "[Access Key]"
	secretKey := "[Secret Key]"

	baseURL := "https://api.hbdm.com"
	//baseURL := "https://api.btcgateway.pro"
	apiParams := &hbdm.ApiParameter{
		Debug:              true,
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		EnablePrivateSign:  false,
		Url:                baseURL,
		PrivateKeyPrime256: "",
	}
	client := hbdm.NewClient(apiParams)

	client.GetAccountInfo("BTC")
	orderResult, err := client.Order(
		"BTC",
		"this_week",
		"",
		0,
		3000.0,
		1,
		"buy",
		"open",
		10,
		"limit",
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("%#v", orderResult)

	orders, err := client.GetOpenOrders(
		"BTC",
		0,
		0,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("%#v", orders)
}
```