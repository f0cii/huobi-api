package hbdmswap

// 接口文档
// https://docs.huobigroup.com/docs/coin_margined_swap/v1/cn/#websocket

import (
	"context"
	"fmt"
	"github.com/frankrap/huobi-api/util"
	"github.com/lithammer/shortuuid/v3"
	"github.com/recws-org/recws"
	"github.com/tidwall/gjson"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

// NWS WebSocket 订单和用户数据接口
type NWS struct {
	sync.RWMutex

	ctx    context.Context
	cancel context.CancelFunc
	conn   recws.RecConn

	wsURL         string
	cid           string // 客户端请求唯一ID
	accessKey     string
	secretKey     string
	subscriptions map[string]interface{}

	ordersCallback            func(order *WSOrder)
	accountsCallback          func(accounts *WSAccounts)
	positionsCallback         func(positions *WSPositions)
	liquidationOrdersCallback func(liquidationOrders *WSLiquidationOrders)
}

// SetProxy 设置代理地址
// porxyURL:
// socks5://127.0.0.1:1080
// https://127.0.0.1:1080
func (ws *NWS) SetProxy(proxyURL string) (err error) {
	var purl *url.URL
	purl, err = url.Parse(proxyURL)
	if err != nil {
		return
	}
	log.Printf("[ws][%s] proxy url:%s", proxyURL, purl)
	ws.conn.Proxy = http.ProxyURL(purl)
	return
}

// SubscribeOrders 订阅订单成交数据
// symbol: BTC-USD
func (ws *NWS) SubscribeOrders(id string, symbol string) {
	ws.SubscribeTopic(id,
		fmt.Sprintf("orders.%s", symbol))
}

// SubscribeAccounts 订阅资产变动数据
// symbol: BTC-USD
func (ws *NWS) SubscribeAccounts(id string, symbol string) {
	ws.SubscribeTopic(id,
		fmt.Sprintf("accounts.%s", symbol))
}

// SubscribePositions 订阅持仓变动更新数据
// symbol: BTC-USD
func (ws *NWS) SubscribePositions(id string, symbol string) {
	ws.SubscribeTopic(id,
		fmt.Sprintf("positions.%s", strings.ToLower(symbol)))
}

// SubscribeLiquidationOrders 订阅强平订单数据
// symbol: BTC-USD
func (ws *NWS) SubscribeLiquidationOrders(id string, symbol string) {
	ws.SubscribeTopic(id,
		fmt.Sprintf("liquidationOrders.%s", strings.ToLower(symbol)))
}

// SubscribeTopic 订阅
func (ws *NWS) SubscribeTopic(id string, topic string) {
	ch := map[string]interface{}{
		"op":    "sub",
		"cid":   ws.cid,
		"topic": topic}
	err := ws.Subscribe(id, ch)
	if err != nil {
		log.Printf("%v", err)
	}
}

// Subscribe 订阅
func (ws *NWS) Subscribe(id string, ch map[string]interface{}) error {
	ws.Lock()
	defer ws.Unlock()

	ws.subscriptions[id] = ch
	return ws.sendWSMessage(ch)
}

// Login 授权接口
func (ws *NWS) Login() error {
	log.Printf("Login")
	if ws.accessKey == "" || ws.secretKey == "" {
		return fmt.Errorf("missing accessKey or secretKey")
	}
	opReq := map[string]string{
		"op":   "auth",
		"type": "api",
	}
	err := ws.setSignatureData(opReq, ws.accessKey, ws.secretKey)
	if err != nil {
		return err
	}
	log.Printf("opReq: %#v", opReq)
	return ws.sendWSMessage(opReq)
}

func (ws *NWS) setSignatureData(data map[string]string, apiKey, apiSecretKey string) error {
	data["AccessKeyId"] = apiKey
	data["SignatureMethod"] = "HmacSHA256"
	data["SignatureVersion"] = "2"
	data["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05")
	postForm := url.Values{}
	// 当type为api时，参数op，type，cid，Signature不参加签名计算
	isApi := data["type"] == "api"
	for k, v := range data {
		if isApi && (k == "op" || k == "cid" || k == "type") {
			continue
		}
		postForm.Set(k, v)
	}
	u, err := url.Parse(ws.wsURL)
	if err != nil {
		return err
	}
	payload := fmt.Sprintf("%s\n%s\n%s\n%s", "GET", u.Host, u.Path, postForm.Encode())
	sign, _ := util.GetParamHmacSHA256Base64Sign(apiSecretKey, payload)
	data["Signature"] = sign
	return nil
}

func (ws *NWS) SetOrdersCallback(callback func(order *WSOrder)) {
	ws.ordersCallback = callback
}

func (ws *NWS) SetAccountsCallback(callback func(accounts *WSAccounts)) {
	ws.accountsCallback = callback
}

func (ws *NWS) SetPositionsCallback(callback func(positions *WSPositions)) {
	ws.positionsCallback = callback
}

func (ws *NWS) SetLiquidationOrdersCallback(callback func(liquidationOrders *WSLiquidationOrders)) {
	ws.liquidationOrdersCallback = callback
}

// Unsubscribe 取消订阅
func (ws *NWS) Unsubscribe(id string) error {
	ws.Lock()
	defer ws.Unlock()

	if v, ok := ws.subscriptions[id]; ok {
		ch, ok := v.(map[string]interface{})
		if ok {
			ch["op"] = "unsub"
			log.Printf("取消订阅: %#v", ch)
			ws.sendWSMessage(ch)
		}
		delete(ws.subscriptions, id)
	}
	return nil
}

func (ws *NWS) unsubscribe(symbol string) {
	ch := map[string]interface{}{
		"op":    "unsub",
		"cid":   ws.cid,
		"topic": fmt.Sprintf("matchOrders.%s", strings.ToLower(symbol))}
	ws.sendWSMessage(ch)
}

func (ws *NWS) subscribeHandler() error {
	//log.Printf("subscribeHandler")
	ws.Lock()
	defer ws.Unlock()

	// 授权
	if ws.accessKey != "" && ws.secretKey != "" {
		err := ws.Login()
		if err != nil {
			log.Printf("%v", err)
			return err
		}
		//time.Sleep(1*time.Second)
	}

	for _, v := range ws.subscriptions {
		log.Printf("sub: %#v", v)
		err := ws.sendWSMessage(v)
		if err != nil {
			log.Printf("%v", err)
		}
	}
	return nil
}

func (ws *NWS) sendWSMessage(msg interface{}) error {
	return ws.conn.WriteJSON(msg)
}

func (ws *NWS) Start() {
	ws.conn.Dial(ws.wsURL, nil)
	go ws.run()
}

func (ws *NWS) run() {
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

func (ws *NWS) handleMsg(messageType int, msg []byte) {
	ret := gjson.ParseBytes(msg)

	if opValue := ret.Get("op"); opValue.Exists() {
		op := opValue.String()
		if op == "ping" {
			ts := ret.Get("ts").Int()
			ws.handlePing(ts)
			return
		} else if op == "notify" {
			topicValue := ret.Get("topic")
			if !topicValue.Exists() {
				log.Printf("err")
				return
			}
			topic := topicValue.String()
			ws.handleNotify(topic, msg...)
		}
	}

	//log.Printf("%v", string(msg))
}

func (ws *NWS) handleNotify(topic string, msg ...byte) {
	if strings.HasPrefix(topic, "orders.") {
		var value WSOrder
		err := json.Unmarshal(msg, &value)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		if ws.ordersCallback != nil {
			ws.ordersCallback(&value)
		}
		return
		//} else if strings.HasPrefix(topic, "matchOrders.") {
		//	var value WSUserMatchOrder
		//	err := json.Unmarshal(msg, &value)
		//	if err != nil {
		//		log.Printf("%v", err)
		//		return
		//	}
		//
		//	if ws.matchOrdersCallback != nil {
		//		ws.matchOrdersCallback(&value)
		//	}
		//	return
	} else if strings.HasPrefix(topic, "accounts.") {
		var value WSAccounts
		err := json.Unmarshal(msg, &value)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		if ws.accountsCallback != nil {
			ws.accountsCallback(&value)
		}
		return
	} else if strings.HasPrefix(topic, "positions.") {
		var value WSPositions
		err := json.Unmarshal(msg, &value)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		if ws.positionsCallback != nil {
			ws.positionsCallback(&value)
		}
		return
	} else if strings.HasPrefix(topic, "liquidationOrders.") {
		var value WSLiquidationOrders
		err := json.Unmarshal(msg, &value)
		if err != nil {
			log.Printf("%v", err)
			return
		}

		if ws.liquidationOrdersCallback != nil {
			ws.liquidationOrdersCallback(&value)
		}
		return
	}
}

func (ws *NWS) handlePing(ts int64) {
	pong := struct {
		Op string `json:"op"`
		Ts int64  `json:"ts"`
	}{"pong", ts}

	err := ws.sendWSMessage(pong)
	if err != nil {
		log.Printf("Send pong error: %v", err)
	}
}

// NewNWS 创建 NWS
// wsURL:
// 正式地址 wss://api.hbdm.com/swap-notification
// 开发地址 wss://api.btcgateway.pro/swap-notification
func NewNWS(wsURL string, accessKey string, secretKey string) *NWS {
	ws := &NWS{
		wsURL:         wsURL,
		cid:           shortuuid.New(),
		accessKey:     accessKey,
		secretKey:     secretKey,
		subscriptions: make(map[string]interface{}),
	}
	ws.ctx, ws.cancel = context.WithCancel(context.Background())
	ws.conn = recws.RecConn{
		KeepAliveTimeout: 10 * time.Second,
	}
	ws.conn.SubscribeHandler = ws.subscribeHandler
	return ws
}
