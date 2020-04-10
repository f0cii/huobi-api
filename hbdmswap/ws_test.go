package hbdmswap

import (
	"log"
	"testing"
)

func newWSTest() *WS {
	wsURL := "wss://api.btcgateway.pro/swap-ws"
	ws := NewWS(wsURL, "", "")
	return ws
}

func TestWS_SubscribeTicker(t *testing.T) {
	ws := newWSTest()
	ws.SetTickerCallback(func(ticker *WSTicker) {
		log.Printf("ticker: %#v", ticker)
	})
	ws.SubscribeTicker("ticker_1", "BTC-USD")
	ws.Start()

	select {}
}

func TestWS_SubscribeDepth(t *testing.T) {
	ws := newWSTest()

	ws.SetDepthCallback(func(depth *WSDepth) {
		log.Printf("depth: %#v", depth)
	})
	ws.SubscribeDepth("depth_1", "BTC-USD")
	ws.Start()

	select {}
}

func TestWS_SubscribeTrade(t *testing.T) {
	ws := newWSTest()

	ws.SetTradeCallback(func(trade *WSTrade) {
		log.Printf("trade: %#v", trade)
	})
	ws.SubscribeTrade("trade_1", "BTC-USD")
	ws.Start()

	select {}
}
