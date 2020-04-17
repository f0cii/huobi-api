package hbdmswap

import (
	"github.com/spf13/viper"
	"log"
	"testing"
)

func newTestClient() *Client {
	viper.SetConfigName("test_config")
	viper.AddConfigPath("..")
	err := viper.ReadInConfig()
	if err != nil {
		log.Panic(err)
	}

	accessKey := viper.GetString("access_key")
	secretKey := viper.GetString("secret_key")

	baseURL := "https://api.btcgateway.pro"
	// baseURL := https://api.hbdm.com
	apiParams := &ApiParameter{
		Debug:              true,
		AccessKey:          accessKey,
		SecretKey:          secretKey,
		EnablePrivateSign:  false,
		Url:                baseURL,
		PrivateKeyPrime256: "",
	}
	c := NewClient(apiParams)
	return c
}

func TestClient_Heartbeat(t *testing.T) {
	c := newTestClient()
	ret, err := c.Heartbeat()
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ret)
}
