package hbdm

import "testing"

func TestClient_Order(t *testing.T) {
	c := newTestClient()
	orderResult, err := c.Order(
		"BTC",
		"this_week",
		"",
		0,
		3000.0,
		1,
		"buy",
		"open",
		5,
		"limit",
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orderResult)
}

func TestClosePosition(t *testing.T) {
	c := newTestClient()
	orderResult, err := c.Order(
		"BTC",
		"this_week",
		"",
		0,
		6000.0, // Buy price must be lower than 6795.94 USD. Sell price must exceed 6294.48 USD.
		1,
		"buy",
		"close",
		5,
		"limit",
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orderResult)
}

func TestClient_Cancel(t *testing.T) {
	c := newTestClient()
	orderID := int64(690494908993323008) // 690495528999559168
	cancelResult, err := c.Cancel("BTC", orderID, 0)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", cancelResult)
}

func TestClient_OrderInfo(t *testing.T) {
	c := newTestClient()
	info, err := c.OrderInfo(
		"BTC",
		696107957142069248,
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", info)
}

func TestClient_GetOpenOrders(t *testing.T) {
	c := newTestClient()
	ordersResult, err := c.GetOpenOrders(
		"BTC",
		0,
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", ordersResult)
}

func TestClient_GetHisOrders(t *testing.T) {
	c := newTestClient()
	ordersResult, err := c.GetHisOrders(
		"BTC",
		0,
		1,
		0,
		90,
		1,
		50,
		"",
		"1",
	)
	if err != nil {
		t.Error(err)
		return
	}
	//t.Logf("%#v", ordersResult)

	for _, v := range ordersResult.Data.Orders {
		t.Logf("%#v", v)
	}
}

func TestClient_LightningClosePosition(t *testing.T) {
	c := newTestClient()
	orderResult, err := c.LightningClosePosition(
		"BTC",
		"this_week",
		"",
		1,
		"sell",
		0,
		"",
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orderResult)
}
