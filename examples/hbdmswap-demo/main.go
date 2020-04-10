package main

import (
	"github.com/frankrap/huobi-api/hbdmswap"
	"log"
)

func main() {
	accessKey := "[Access Key]"
	secretKey := "[Secret Key]"

	baseURL := "https://api.hbdm.com"
	//baseURL := "https://api.btcgateway.pro"
	apiParams := &hbdmswap.ApiParameter{
		Debug:              true,
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		EnablePrivateSign:  false,
		Url:                baseURL,
		PrivateKeyPrime256: "",
	}
	client := hbdmswap.NewClient(apiParams)

	symbol := "BTC-USD"
	client.GetAccountInfo(symbol)
	orderResult, err := client.Order(symbol,
		0,
		3000,
		1,
		"buy",
		"open",
		125,
		"limit")
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("%#v", orderResult)

	orders, err := client.GetOpenOrders(
		symbol,
		0,
		0,
	)
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Printf("%#v", orders)
}
