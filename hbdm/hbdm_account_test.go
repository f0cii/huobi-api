package hbdm

import "testing"

func TestClient_GetAccountInfo(t *testing.T) {
	c := newTestClient()
	info, err := c.GetAccountInfo("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", info)
}

func TestClient_GetPositionInfo(t *testing.T) {
	c := newTestClient()
	info, err := c.GetPositionInfo("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%v", info)
}
