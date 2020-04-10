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
	orderResult, err := client.Order("BTC",
		"this_week",
		"",
		0,
		3000.0,
		1,
		"buy",
		"open",
		10,
		"limit")
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
