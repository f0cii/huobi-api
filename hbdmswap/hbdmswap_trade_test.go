package hbdmswap

import "testing"

func TestClient_Order(t *testing.T) {
	c := newTestClient()
	orderResult, err := c.Order(
		"BTC-USD",
		0,
		3000,
		1,
		"buy",
		"open",
		125,
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
	r, err := c.Cancel(
		"BTC-USD",
		696382271304871936,
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", r)
}

func TestClient_OrderInfo(t *testing.T) {
	c := newTestClient()
	info, err := c.OrderInfo(
		"BTC-USD",
		696382271304871936,
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
	orders, err := c.GetOpenOrders(
		"BTC-USD",
		1,
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%#v", orders)
}

func TestClient_GetHisOrders(t *testing.T) {
	c := newTestClient()
	orders, err := c.GetHisOrders(
		"BTC-USD",
		0,
		1,
		0,
		90,
		1,
		0,
	)
	if err != nil {
		t.Error(err)
		return
	}

	for _, v := range orders.Data.Orders {
		t.Logf("%#v", v)
	}
}
