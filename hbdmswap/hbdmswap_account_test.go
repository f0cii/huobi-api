package hbdmswap

import "testing"

func TestClient_GetAccountInfo(t *testing.T) {
	c := newTestClient()
	accountInfo, err := c.GetAccountInfo("BTC-USD")
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range accountInfo.Data {
		t.Logf("%#v", v)
	}
}

func TestClient_GetPositionInfo(t *testing.T) {
	c := newTestClient()
	positions, err := c.GetPositionInfo("BTC-USD")
	if err != nil {
		t.Error(err)
		return
	}
	for _, v := range positions.Data {
		t.Logf("%#v", v)
	}
}
