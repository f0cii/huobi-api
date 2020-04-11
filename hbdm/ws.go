package hbdm

// 接口文档
// https://huobiapi.github.io/docs/dm/v1/cn/#websocket

import (
	"context"
	"fmt"
	"github.com/frankrap/huobi-api/util"
	"github.com/google/uuid"
	"github.com/recws-org/recws"
	"github.com/tidwall/gjson"
	"log"
	"strings"
	"sync"
	"time"
)

// WS WebSocket 市场行情接口
type WS struct {
	sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	wsConn recws.RecConn

	wsURL         string // Public
	cid           string // 客户端请求唯一ID
	accessKey     string
	secretKey     string
	subscriptions map[string]interface{}

	tickerCallback func(trade *WSTicker)
	depthCallback  func(depth *WSDepth)
	tradeCallback  func(trade *WSTrade)
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
	return ws.sendWSMessage(ch)
}

func (ws *WS) SetTickerCallback(callback func(ticker *WSTicker)) {
	ws.tickerCallback = callback
}

func (ws *WS) SetDepthCallback(callback func(depth *WSDepth)) {
	ws.depthCallback = callback
}

func (ws *WS) SetTradeCallback(callback func(trade *WSTrade)) {
	ws.tradeCallback = callback
}

// Unsubscribe 取消订阅
func (ws *WS) Unsubscribe(id string) error {
	ws.Lock()
	defer ws.Unlock()

	if _, ok := ws.subscriptions[id]; ok {
		delete(ws.subscriptions, id)
	}
	return nil
}

func (ws *WS) subscribeHandler() error {
	//log.Printf("subscribeHandler")
	ws.Lock()
	defer ws.Unlock()

	for _, v := range ws.subscriptions {
		//log.Printf("sub: %#v", v)
		err := ws.sendWSMessage(v)
		if err != nil {
			log.Printf("%v", err)
		}
	}
	return nil
}

func (ws *WS) sendWSMessage(msg interface{}) error {
	return ws.wsConn.WriteJSON(msg)
}

func (ws *WS) Start() {
	ws.wsConn.Dial(ws.wsURL, nil)
	go ws.run()
}

func (ws *WS) run() {
	ctx := context.Background()
	for {
		select {
		case <-ctx.Done():
			go ws.wsConn.Close()
			log.Printf("Websocket closed %s", ws.wsConn.GetURL())
			return
		default:
			messageType, msg, err := ws.wsConn.ReadMessage()
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

	fmt.Printf("%v", string(msg))

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
			if strings.Contains(ch, ".depth") {
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
		} else {
			log.Printf("msg: %v", string(msg))
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
// 正式地址 wss://api.hbdm.com/ws
// 开发地址 wss://api.btcgateway.pro/ws
func NewWS(wsURL string, accessKey string, secretKey string) *WS {
	ws := &WS{
		wsURL:         wsURL,
		cid:           uuid.New().String(),
		accessKey:     accessKey,
		secretKey:     secretKey,
		subscriptions: make(map[string]interface{}),
	}
	ws.ctx, ws.cancel = context.WithCancel(context.Background())
	ws.wsConn = recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	ws.wsConn.SubscribeHandler = ws.subscribeHandler
	return ws
}
