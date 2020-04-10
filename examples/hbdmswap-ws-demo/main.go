package main

import (
	"github.com/frankrap/huobi-api/hbdmswap"
	"log"
)

func main() {
	//wsURL := "wss://api.hbdm.com/swap-ws"
	wsURL := "wss://api.btcgateway.pro/swap-ws"
	ws := hbdmswap.NewWS(wsURL, "", "")

	// 设置Ticker回调
	ws.SetTickerCallback(func(ticker *hbdmswap.WSTicker) {
		log.Printf("ticker: %#v", ticker)
	})
	// 设置Depth回调
	ws.SetDepthCallback(func(depth *hbdmswap.WSDepth) {
		log.Printf("depth: %#v", depth)
	})
	// 设置Trade回调
	ws.SetTradeCallback(func(trade *hbdmswap.WSTrade) {
		log.Printf("trade: %#v", trade)
	})

	// 订阅Ticker
	ws.SubscribeTicker("ticker_1", "BTC-USD")
	// 订阅Depth
	ws.SubscribeDepth("depth_1", "BTC-USD")
	// 订阅Trade
	ws.SubscribeTrade("trade_1", "BTC-USD")
	// 启动WS
	ws.Start()

	select {}
}
