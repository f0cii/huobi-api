package hbdm

import "testing"

func TestClient_ContractInfo(t *testing.T) {
	c := newTestClient()
	contractInfo, err := c.GetContractInfo(
		"BTC",
		"this_week",
		"",
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", contractInfo)
}

func TestClient_GetContractIndex(t *testing.T) {
	c := newTestClient()
	info, err := c.GetContractIndex("BTC")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", info)
}

func TestClient_GetMarketDepth(t *testing.T) {
	c := newTestClient()
	depth, err := c.GetMarketDepth("BTC_CQ",
		"step5")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", depth)
}

func TestClient_GetKLine(t *testing.T) {
	c := newTestClient()
	klineResult, err := c.GetKLine(
		"BTC_CQ",
		"1min",
		0,
		0,
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", klineResult)
}
