package cc_client

import (
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
	ApiUrl        string
	WebSocketUrl  string
	ApiKey        string
	ApiSecret     string
	HTTPClient    *http.Client
	WebSocketConn *websocket.Conn
}

type RequestOptions struct {
	Authorization bool
	Params        url.Values
}

func NewClient(apiKey string, apiSecret string) *Client {
	return &Client{
		ApiUrl:       ApiUrl,
		WebSocketUrl: WebSocketUrl,
		ApiKey:       apiKey,
		ApiSecret:    apiSecret,
		HTTPClient: &http.Client{
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

func (c *Client) SubscribeOrderBooks(listener func(pair string, diff *OrderBooks)) error {
	if err := c.connectWebSocket(); err != nil {
		return err
	}

	go func() {
		for {
			// TODO 途中でエラーになったときにgoroutine再実行できるようにしたい
			_, message, err := c.WebSocketConn.ReadMessage()
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

	return c.WebSocketConn.WriteMessage(websocket.TextMessage, message)
}

func (c *Client) connectWebSocket() error {
	if c.WebSocketConn != nil {
		return nil
	}

	conn, _, err := websocket.DefaultDialer.Dial(c.WebSocketUrl, nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	c.WebSocketConn = conn

	return nil
}

func (c *Client) sendRequest(method string, path string, options *RequestOptions, responseBody interface{}) error {
	url := c.ApiUrl + path
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	if options.Authorization {
		nonce := strconv.FormatInt(time.Now().UnixNano(), 10)
		body := ""
		signature := util.HmacSha256(nonce+url+body, c.ApiSecret)

		req.Header.Set("ACCESS-KEY", c.ApiKey)
		req.Header.Set("ACCESS-NONCE", nonce)
		req.Header.Set("ACCESS-SIGNATURE", signature)
	}
	if options.Params != nil {
		req.URL.RawQuery = options.Params.Encode()
	}

	res, err := c.HTTPClient.Do(req)
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
		fmt.Println(err)
		return err
	}

	return nil
}
