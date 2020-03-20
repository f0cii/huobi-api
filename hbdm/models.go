package hbdm

import "encoding/json"

type ContractInfo struct {
	Symbol         string  `json:"symbol"`
	ContractCode   string  `json:"contract_code"`
	ContractType   string  `json:"contract_type"`
	ContractSize   float64 `json:"contract_size"`
	PriceTick      float64 `json:"price_tick"`
	DeliveryDate   string  `json:"delivery_date"`
	CreateDate     string  `json:"create_date"`
	ContractStatus int     `json:"contract_status"`
}

type ContractInfoResult struct {
	Status string         `json:"status"`
	Data   []ContractInfo `json:"data"`
	Ts     int64          `json:"ts"`
}

type ContractIndex struct {
	Symbol     string  `json:"symbol"`
	IndexPrice float64 `json:"index_price"`
	IndexTs    int64   `json:"index_ts"`
}

type ContractIndexResult struct {
	Status string          `json:"status"`
	Data   []ContractIndex `json:"data"`
	Ts     int64           `json:"ts"`
}

type Tick struct {
	Asks    [][]float64 `json:"asks"`
	Bids    [][]float64 `json:"bids"`
	Ch      string      `json:"ch"`
	ID      int         `json:"id"`
	MrID    int64       `json:"mrid"`
	Ts      int64       `json:"ts"`
	Version int         `json:"version"`
}

type MarketDepthResult struct {
	Ch     string `json:"ch"`
	Status string `json:"status"`
	Tick   Tick   `json:"tick"`
	Ts     int64  `json:"ts"`
}

type KLine struct {
	Amount float64 `json:"amount"`
	Close  float64 `json:"close"`
	Count  int     `json:"count"`
	High   float64 `json:"high"`
	ID     int     `json:"id"`
	Low    float64 `json:"low"`
	Open   float64 `json:"open"`
	Vol    int     `json:"vol"`
}

type KLineResult struct {
	Ch     string  `json:"ch"`
	Data   []KLine `json:"data"`
	Status string  `json:"status"`
	Ts     int64   `json:"ts"`
}

type AccountInfo struct {
	Symbol            string      `json:"symbol"`
	MarginBalance     float64     `json:"margin_balance"`
	MarginPosition    float64     `json:"margin_position"`
	MarginFrozen      float64     `json:"margin_frozen"`
	MarginAvailable   float64     `json:"margin_available"`
	ProfitReal        float64     `json:"profit_real"`
	ProfitUnreal      float64     `json:"profit_unreal"`
	RiskRate          interface{} `json:"risk_rate"`
	WithdrawAvailable float64     `json:"withdraw_available"`
	LiquidationPrice  interface{} `json:"liquidation_price"`
	LeverRate         float64     `json:"lever_rate"`
	AdjustFactor      float64     `json:"adjust_factor"`
	MarginStatic      float64     `json:"margin_static"`
	IsDebit           int         `json:"is_debit"`
}

type AccountInfoResult struct {
	Status string        `json:"status"`
	Data   []AccountInfo `json:"data"`
	Ts     int64         `json:"ts"`
}

type PositionInfoResult struct {
	Status string        `json:"status"`
	Data   []interface{} `json:"data"`
	Ts     int64         `json:"ts"`
}

type OrderData struct {
	OrderID       int64  `json:"order_id"`
	OrderIDStr    string `json:"order_id_str"`
	ClientOrderID int64  `json:"client_order_id"`
}

type OrderResult struct {
	Status string    `json:"status"`
	Data   OrderData `json:"data"`
	Ts     int64     `json:"ts"`
}

type CancelError struct {
	OrderID string `json:"order_id"`
	ErrCode int    `json:"err_code"`
	ErrMsg  string `json:"err_msg"`
}

type CancelData struct {
	Errors    []CancelError `json:"errors"`
	Successes string        `json:"successes"`
}

type CancelResult struct {
	Status string     `json:"status"`
	Data   CancelData `json:"data"`
	Ts     int64      `json:"ts"`
}

type Order struct {
	Symbol          string      `json:"symbol"`
	ContractCode    string      `json:"contract_code"`
	ContractType    string      `json:"contract_type"`
	Volume          int         `json:"volume"`
	Price           float64     `json:"price"`
	OrderPriceType  json.Number `json:"order_price_type"`
	OrderType       int         `json:"order_type"`
	Direction       string      `json:"direction"`
	Offset          string      `json:"offset"`
	LeverRate       int         `json:"lever_rate"`
	OrderID         int64       `json:"order_id"`
	ClientOrderID   string      `json:"client_order_id"`
	CreatedAt       int64       `json:"created_at"`
	TradeVolume     float64     `json:"trade_volume"`
	TradeTurnover   float64     `json:"trade_turnover"`
	Fee             float64     `json:"fee"`
	TradeAvgPrice   float64     `json:"trade_avg_price"`
	MarginFrozen    float64     `json:"margin_frozen"`
	Profit          float64     `json:"profit"`
	Status          int         `json:"status"`
	OrderSource     string      `json:"order_source"`
	OrderIDStr      string      `json:"order_id_str"`
	FeeAsset        string      `json:"fee_asset"`
	LiquidationType string      `json:"liquidation_type"`
	CreateDate      int64       `json:"create_date"`
}

type OrderInfoResult struct {
	Status string  `json:"status"`
	Data   []Order `json:"data"`
	Ts     int64   `json:"ts"`
}

type OpenOrdersData struct {
	Orders      []Order `json:"orders"`
	TotalPage   int     `json:"total_page"`
	CurrentPage int     `json:"current_page"`
	TotalSize   int     `json:"total_size"`
}

type OpenOrdersResult struct {
	Status string         `json:"status"`
	Data   OpenOrdersData `json:"data"`
	Ts     int64          `json:"ts"`
}

type HisOrdersData struct {
	Orders      []Order `json:"orders"`
	TotalPage   int     `json:"total_page"`
	CurrentPage int     `json:"current_page"`
	TotalSize   int     `json:"total_size"`
}

type HisOrdersResult struct {
	Status string        `json:"status"`
	Data   HisOrdersData `json:"data"`
	Ts     int64         `json:"ts"`
}

type LightningClosePositionResult struct {
	Status string    `json:"status"`
	Data   OrderData `json:"data"`
	Ts     int64     `json:"ts"`
}
