package hbdmswap

import "testing"

func TestClient_GetMarketDepth(t *testing.T) {
	b := newTestClient()
	depth, err := b.GetMarketDepth("BTC-USD", "step6")
	if err != nil {
		t.Error(err)
		return
	}

	//t.Logf("%#v", depth)

	for _, v := range depth.Tick.Asks {
		t.Logf("%#v", v)
	}
}

func TestClient_GetKLine(t *testing.T) {
	b := newTestClient()
	kLine, err := b.GetKLine("BTC-USD", "1min", 10, 0, 0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", kLine)
}
