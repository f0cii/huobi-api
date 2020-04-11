package hbdmswap

import (
	"github.com/spf13/viper"
	"log"
	"testing"
)

func newNWSTest() *NWS {
	viper.SetConfigName("test_config")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")

	wsURL := "wss://api.btcgateway.pro/swap-notification"
	ws := NewNWS(wsURL, accessKey, secretKey)
	return ws
}

func TestNWS_SubscribeOrders(t *testing.T) {
	ws := newNWSTest()

	// 设置回调
	ws.SetOrdersCallback(func(order *WSOrder) {
		log.Printf("order: %#v", order)
	})

	// 订阅
	ws.SubscribeOrders("orders_1", "BTC-USD")
	// 取消订阅
	//ws.Unsubscribe("orders_1")

	//go func() {
	//	time.Sleep(30*time.Second)
	//	ws.Unsubscribe("orders_1")
	//}()

	ws.Start()

	select {}
}
