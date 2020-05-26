package spot

type Tick struct {
	Bids [][]float64 `json:"bids"` // [8956.46,0.225618]
	Asks [][]float64 `json:"asks"`
}

type MarketDepthResult struct {
	Ch     string `json:"ch"`     // market.btcusdt.depth.step0
	Status string `json:"status"` // ok
	Ts     int64  `json:"ts"`     // 1590479896449
	Tick   Tick   `json:"tick"`
}
