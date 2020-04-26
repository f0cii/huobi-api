package hbdmswap

// 接口文档
// https://huobiapi.github.io/docs/coin_margined_swap/v1/cn/#websocket

import (
	"context"
	"fmt"
	"github.com/frankrap/huobi-api/util"
	"github.com/recws-org/recws"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// WS WebSocket 市场行情接口
type WS struct {
	sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	conn   recws.RecConn

	wsURL         string
	accessKey     string
	secretKey     string
	debugMode     bool
	subscriptions map[string]interface{}

	tickerCallback  func(trade *WSTicker)
	depthCallback   func(depth *WSDepth)
	depthHFCallback func(depth *WSDepthHF)
	tradeCallback   func(trade *WSTrade)
}

// SetProxy 设置代理地址
// porxyURL:
// socks5://127.0.0.1:1080
// https://127.0.0.1:1080
func (ws *WS) SetProxy(proxyURL string) (err error) {
	var purl *url.URL
	purl, err = url.Parse(proxyURL)
	if err != nil {
		return
	}
	log.Printf("[ws][%s] proxy url:%s", proxyURL, purl)
	ws.conn.Proxy = http.ProxyURL(purl)
	return
}

// SubscribeTicker 订阅 Market Ticker 数据
// id: 订阅的编号
// symbol: BTC_CQ
func (ws *WS) SubscribeTicker(id string, symbol string) {
	ch := map[string]interface{}{
		"id":  id,
		"sub": fmt.Sprintf("market.%s.detail", symbol)}
	ws.Subscribe(id, ch)
}

// SubscribeDepth 订阅 Market Depth 数据
// id: 订阅的编号
// symbol: BTC_CQ
func (ws *WS) SubscribeDepth(id string, symbol string) {
	ch := map[string]interface{}{
		"id":  id,
		"sub": fmt.Sprintf("market.%s.depth.step0", symbol)}
	ws.Subscribe(id, ch)
}

// SubscribeDepthHF 订阅增量深度
// size: 20/150 档位数，20:表示20档不合并的深度，150:表示150档不合并的深度
// dateType: 数据类型，不填默认为全量数据，"incremental"：增量数据，"snapshot"：全量数据
func (ws *WS) SubscribeDepthHF(id string, symbol string, size int, dateType string) {
	ch := map[string]interface{}{
		"id":        id,
		"sub":       fmt.Sprintf("market.%v.depth.size_%v.high_freq", symbol, size),
		"data_type": dateType,
	}
	ws.Subscribe(id, ch)
}

// SubscribeTrade 订阅 Market Trade 数据
// id: 订阅的编号
// symbol: BTC_CQ
func (ws *WS) SubscribeTrade(id string, symbol string) {
	ch := map[string]interface{}{
		"id":  id,
		"sub": fmt.Sprintf("market.%s.trade.detail", symbol)}
	ws.Subscribe(id, ch)
}

// Subscribe 订阅
func (ws *WS) Subscribe(id string, ch map[string]interface{}) error {
	ws.Lock()
	defer ws.Unlock()

	ws.subscriptions[id] = ch
	ws.sendWSMessage(ch)
	return nil
}

func (ws *WS) SetTickerCallback(callback func(ticker *WSTicker)) {
	ws.tickerCallback = callback
}

func (ws *WS) SetDepthCallback(callback func(depth *WSDepth)) {
	ws.depthCallback = callback
}

func (ws *WS) SetDepthHFCallback(callback func(depth *WSDepthHF)) {
	ws.depthHFCallback = callback
}

func (ws *WS) SetTradeCallback(callback func(trade *WSTrade)) {
	ws.tradeCallback = callback
}

func (ws *WS) Unsubscribe(id string) error {
	ws.Lock()
	defer ws.Unlock()

	if _, ok := ws.subscriptions[id]; ok {
		delete(ws.subscriptions, id)
	}
	return nil
}

func (ws *WS) subscribeHandler() error {
	log.Printf("subscribeHandler")
	ws.Lock()
	defer ws.Unlock()

	for _, v := range ws.subscriptions {
		ws.sendWSMessage(v)
	}
	return nil
}

func (ws *WS) sendWSMessage(msg interface{}) error {
	return ws.conn.WriteJSON(msg)
}

func (ws *WS) Start() {
	ws.conn.Dial(ws.wsURL, nil)
	go ws.run()
}

func (ws *WS) run() {
	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			go ws.conn.Close()
			log.Printf("Websocket closed %s", ws.conn.GetURL())
			return
		default:
			messageType, msg, err := ws.conn.ReadMessage()
			if err != nil {
				log.Printf("Read error: %v", err)
				time.Sleep(100 * time.Millisecond)
				continue
			}

			msg, err = util.GzipUncompress(msg)
			if err != nil {
				continue
			}

			ws.handleMsg(messageType, msg)
		}
	}
}

func (ws *WS) handleMsg(messageType int, msg []byte) {
	ret := gjson.ParseBytes(msg)

	if ws.debugMode {
		log.Printf("%v", string(msg))
	}

	if pingValue := ret.Get("ping"); pingValue.Exists() {
		// 心跳
		ping := pingValue.Int()
		ws.handlePing(ping)
		return
	}

	// 订阅成功返回消息
	// {"id":"depth_1","subbed":"market.BTC_CQ.depth.step0","ts":1586498957314,"status":"ok"}

	if chValue := ret.Get("ch"); chValue.Exists() {
		// market.BTC_CQ.depth.step0
		ch := chValue.String()
		if strings.HasPrefix(ch, "market") {
			if strings.HasSuffix(ch, ".high_freq") {
				var depth WSDepthHF
				err := json.Unmarshal(msg, &depth)
				if err != nil {
					log.Printf("%v", err)
					return
				}

				if ws.depthHFCallback != nil {
					ws.depthHFCallback(&depth)
				}
			} else if strings.Contains(ch, ".depth") {
				var depth WSDepth
				err := json.Unmarshal(msg, &depth)
				if err != nil {
					log.Printf("%v", err)
					return
				}

				if ws.depthCallback != nil {
					ws.depthCallback(&depth)
				}
			} else if strings.HasSuffix(ch, ".trade.detail") {
				//log.Printf("%v", string(msg))
				var trade WSTrade
				err := json.Unmarshal(msg, &trade)
				if err != nil {
					log.Printf("%v", err)
					return
				}

				if ws.tradeCallback != nil {
					ws.tradeCallback(&trade)
				}
			} else if strings.HasSuffix(ch, ".detail") {
				var ticker WSTicker
				err := json.Unmarshal(msg, &ticker)
				if err != nil {
					log.Printf("%v", err)
					return
				}

				if ws.tickerCallback != nil {
					ws.tickerCallback(&ticker)
				}
			} else {
				log.Printf("%v", string(msg))
			}
		}
		return
	}
}

func (ws *WS) handlePing(ping int64) {
	pong := struct {
		Pong int64 `json:"pong"`
	}{ping}

	err := ws.sendWSMessage(pong)
	if err != nil {
		log.Printf("Send pong error: %v", err)
	}
}

// NewWS 创建 WS
// wsURL:
// 正式地址 wss://api.hbdm.com/swap-ws
// 开发地址 wss://api.btcgateway.pro/swap-ws
func NewWS(wsURL string, accessKey string, secretKey string, debugMode bool) *WS {
	ws := &WS{
		wsURL:         wsURL,
		accessKey:     accessKey,
		secretKey:     secretKey,
		debugMode:     debugMode,
		subscriptions: make(map[string]interface{}),
	}
	ws.ctx, ws.cancel = context.WithCancel(context.Background())
	ws.conn = recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	ws.conn.SubscribeHandler = ws.subscribeHandler
	return ws
}
