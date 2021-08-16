package cc_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"hamster/lib/clog"
	"hamster/lib/util"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	ApiUrl       = "https://coincheck.com/api"
	WebSocketUrl = "wss://ws-api.coincheck.com"
)

type Client struct {
	apiUrl           string
	webSocketUrl     string
	apiKey           string
	apiSecret        string
	httpClient       *http.Client
	orderBooksSocket *websocket.Conn
	tradesSocket     *websocket.Conn
}

type RequestBody map[string]string
type RequestOptions struct {
	Authorization bool
	Params        url.Values
	Body          RequestBody
}

func NewClient(apiKey string, apiSecret string) *Client {
	return &Client{
		apiUrl:       ApiUrl,
		webSocketUrl: WebSocketUrl,
		apiKey:       apiKey,
		apiSecret:    apiSecret,
		httpClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *Client) GetOrderBooks() (*OrderBooks, error) {
	opts := &RequestOptions{
		Authorization: false,
	}
	res := OrderBooks{}
	if err := c.sendRequest("GET", "/order_books", opts, &res); err != nil {
		clog.Logger.Error(err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetRate(orderType string, pair string) (*Rate, error) {
	opts := &RequestOptions{
		Authorization: false,
		Params: url.Values{
			"order_type": {orderType},
			"pair":       {pair},
			"amount":     {"1"},
		},
	}
	res := Rate{}
	if err := c.sendRequest("GET", "/exchange/orders/rate", opts, &res); err != nil {
		clog.Logger.Error(err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetOpenOrders() (*OpenOrders, error) {
	opts := &RequestOptions{
		Authorization: true,
	}
	res := OpenOrders{}
	if err := c.sendRequest("GET", "/exchange/orders/opens", opts, &res); err != nil {
		clog.Logger.Error(err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetRecentTransactions() (*Transactions, error) {
	opts := &RequestOptions{
		Authorization: true,
	}
	res := Transactions{}
	if err := c.sendRequest("GET", "/exchange/orders/transactions", opts, &res); err != nil {
		clog.Logger.Error(err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) GetBalance() (*Balance, error) {
	opts := &RequestOptions{
		Authorization: true,
	}
	res := Balance{}
	if err := c.sendRequest("GET", "/accounts/balance", opts, &res); err != nil {
		clog.Logger.Error(err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) CreateOrder(
	pair string,
	orderType string,
	rate *float64,
	amount *float64,
	marketBuyAmount *float64,
	stopLossRate *float64,
) (*Order, error) {
	body := RequestBody{
		"pair":       pair,
		"order_type": orderType,
	}
	if rate != nil {
		body["rate"] = strconv.FormatFloat(*rate, 'f', 0, 64)
	}
	if amount != nil {
		body["amount"] = strconv.FormatFloat(*amount, 'f', -1, 64)
	}
	if marketBuyAmount != nil {
		body["market_buy_amount"] = strconv.FormatFloat(*marketBuyAmount, 'f', -1, 64)
	}
	if stopLossRate != nil {
		body["stop_loss_rate"] = strconv.FormatFloat(*stopLossRate, 'f', -1, 64)
	}
	opts := &RequestOptions{
		Authorization: true,
		Body:          body,
	}
	res := Order{}
	if err := c.sendRequest("POST", "/exchange/orders", opts, &res); err != nil {
		clog.Logger.Error(err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) CancelOrder(id uint64) error {
	opts := &RequestOptions{
		Authorization: true,
	}
	res := OrderCancel{}
	if err := c.sendRequest("DELETE", fmt.Sprintf("/exchange/orders/%d", id), opts, &res); err != nil {
		clog.Logger.Error(err)
		return err
	}

	return nil

}

func (c *Client) SubscribeOrderBooks(listener func(pair string, diff *OrderBooks)) error {
	if c.orderBooksSocket != nil {
		if err := c.orderBooksSocket.Close(); err != nil {
			return err
		}
		c.orderBooksSocket = nil
	}

	conn, err := c.createWebSocketConnection()
	if err != nil {
		return err
	}
	c.orderBooksSocket = conn

	go func() {
		for {
			// TODO 途中でエラーになったときにgoroutine再実行できるようにしたい
			_, message, err := c.orderBooksSocket.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			stream := &OrderBooksDiffStream{}
			if err = json.Unmarshal(message, &stream); err != nil {
				fmt.Println(err)
				return
			}

			var pair string
			if err = json.Unmarshal(stream[0], &pair); err != nil {
				fmt.Println(err)
				return
			}

			diff := &OrderBooks{}
			if err = json.Unmarshal(stream[1], &diff); err != nil {
				fmt.Println(err)
				return
			}

			listener(pair, diff)
		}
	}()

	message, err := json.Marshal(map[string]string{
		"type":    "subscribe",
		"channel": "btc_jpy-orderbook",
	})
	if err != nil {
		return err
	}

	return c.orderBooksSocket.WriteMessage(websocket.TextMessage, message)
}

func (c *Client) SubscribeTrades(listener func(trade *Trade)) error {
	if c.tradesSocket != nil {
		if err := c.tradesSocket.Close(); err != nil {
			return err
		}
		c.tradesSocket = nil
	}

	conn, err := c.createWebSocketConnection()
	if err != nil {
		return err
	}
	c.tradesSocket = conn

	go func() {
		for {
			// TODO 途中でエラーになったときにgoroutine再実行できるようにしたい
			_, message, err := c.tradesSocket.ReadMessage()
			if err != nil {
				log.Println("read error:", err)
				return
			}

			stream := &TradeStream{}
			if err = json.Unmarshal(message, &stream); err != nil {
				fmt.Println(err)
				return
			}

			trade := &Trade{}
			if err = json.Unmarshal(stream[0], &trade.Id); err != nil {
				fmt.Println(err)
				return
			}
			if err = json.Unmarshal(stream[1], &trade.Pair); err != nil {
				fmt.Println(err)
				return
			}
			if err = json.Unmarshal(stream[2], &trade.Rate); err != nil {
				fmt.Println(err)
				return
			}
			if err = json.Unmarshal(stream[3], &trade.Amount); err != nil {
				fmt.Println(err)
				return
			}
			if err = json.Unmarshal(stream[4], &trade.OrderType); err != nil {
				fmt.Println(err)
				return
			}

			listener(trade)
		}
	}()

	message, err := json.Marshal(map[string]string{
		"type":    "subscribe",
		"channel": "btc_jpy-trades",
	})
	if err != nil {
		return err
	}

	return c.tradesSocket.WriteMessage(websocket.TextMessage, message)
}

func (c *Client) createWebSocketConnection() (*websocket.Conn, error) {
	conn, _, err := websocket.DefaultDialer.Dial(c.webSocketUrl, nil)
	if err != nil {
		clog.Logger.Errorf("failed to connect websocket: %s", err)
		return nil, err
	}
	return conn, nil
}

func (c *Client) sendRequest(method string, path string, options *RequestOptions, responseBody interface{}) (err error) {
	var req *http.Request
	var requestBody string

	requestUrl := c.apiUrl + path

	if method == "POST" && options.Body != nil {
		jsonBytes, err := json.Marshal(options.Body)
		if err != nil {
			return err
		}
		req, err = http.NewRequest(method, requestUrl, bytes.NewBuffer(jsonBytes))
		requestBody = string(jsonBytes)
	} else {
		req, err = http.NewRequest(method, requestUrl, nil)
	}
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	if options.Authorization {
		nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
		signature := util.HmacSha256(nonce+requestUrl+requestBody, c.apiSecret)

		req.Header.Set("ACCESS-KEY", c.apiKey)
		req.Header.Set("ACCESS-NONCE", nonce)
		req.Header.Set("ACCESS-SIGNATURE", signature)
	}
	if options.Params != nil {
		req.URL.RawQuery = options.Params.Encode()
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		bodyBytes, err := ioutil.ReadAll(res.Body)
		if err != nil {
			clog.Logger.Fatal(err)
		}
		bodyString := string(bodyBytes)
		return fmt.Errorf("api request error [%d]: %s", res.StatusCode, bodyString)
	}

	if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		clog.Logger.Fatal(err)
		return err
	}

	return nil
}
