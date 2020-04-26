package hbdm

import (
	"github.com/spf13/viper"
	"log"
	"testing"
)

func newWSTest() *WS {
	viper.SetConfigName("test_config")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")

	wsURL := "wss://api.btcgateway.pro/ws"
	ws := NewWS(wsURL, accessKey, secretKey, true)
	return ws
}

func TestWS_SubscribeTicker(t *testing.T) {
	ws := newWSTest()

	ws.SetTickerCallback(func(ticker *WSTicker) {
		log.Printf("ticker: %#v", ticker)
	})
	ws.SubscribeTicker("ticker_1", "BTC_CQ")
	ws.Start()

	select {}
}

func TestWS_SubscribeDepth(t *testing.T) {
	ws := newWSTest()

	ws.SetDepthCallback(func(depth *WSDepth) {
		log.Printf("depth: %#v", depth)
	})
	ws.SubscribeDepth("depth_1", "BTC_CQ")
	ws.Start()

	select {}
}

func TestWS_SubscribeDepthHighFreq(t *testing.T) {
	ws := newWSTest()

	ws.SetDepthHFCallback(func(depth *WSDepthHF) {
		log.Printf("depth: %#v", depth)
	})
	ws.SubscribeDepthHF("depth_1", "BTC_CQ", 20, "incremental")
	ws.Start()

	select {}
}

func TestWS_SubscribeTrade(t *testing.T) {
	ws := newWSTest()

	ws.SetTradeCallback(func(trade *WSTrade) {
		log.Printf("trade: %#v", trade)
	})
	ws.SubscribeTrade("trade_1", "BTC_CQ")
	ws.Start()

	select {}
}
