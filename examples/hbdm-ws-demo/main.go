package main

import (
	"github.com/frankrap/huobi-api/hbdm"
	"log"
)

func main() {
	wsURL := "wss://api.hbdm.com/ws"
	//wsURL := "wss://api.btcgateway.pro/ws"
	ws := hbdm.NewWS(wsURL, "", "")

	// 设置Ticker回调
	ws.SetTickerCallback(func(ticker *hbdm.WSTicker) {
		log.Printf("ticker: %#v", ticker)
	})
	// 设置Depth回调
	ws.SetDepthCallback(func(depth *hbdm.WSDepth) {
		log.Printf("depth: %#v", depth)
	})
	// 设置Trade回调
	ws.SetTradeCallback(func(trade *hbdm.WSTrade) {
		log.Printf("trade: %#v", trade)
	})

	// 订阅Ticker
	ws.SubscribeTicker("ticker_1", "BTC_CQ")
	// 订阅Depth
	ws.SubscribeDepth("depth_1", "BTC_CQ")
	// 订阅Trade
	ws.SubscribeTrade("trade_1", "BTC_CQ")
	// 启动WS
	ws.Start()

	select {}
}
