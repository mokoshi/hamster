package cc_client

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"time"
)

const (
	ApiUrl       = "https://coincheck.com/api"
	WebSocketUrl = "wss://ws-api.coincheck.com"
)

type Client struct {
	ApiUrl        string
	WebSocketUrl  string
	apiKey        string
	HTTPClient    *http.Client
	WebSocketConn *websocket.Conn
}

func NewClient(apiKey string) *Client {
	return &Client{
		ApiUrl:       ApiUrl,
		WebSocketUrl: WebSocketUrl,
		apiKey:       apiKey,
		HTTPClient: &http.Client{
			Timeout: time.Second * 5,
		},
	}
}

func (c *Client) GetOrderBooks() (*OrderBooks, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/order_books", c.ApiUrl), nil)
	if err != nil {
		return nil, err
	}

	res := OrderBooks{}
	if err := c.sendRequest(req, &res); err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &res, nil
}

func (c *Client) SubscribeOrderBooks(listener func(pair string, diff *OrderBooks)) error {
	if err := c.ConnectWebSocket(); err != nil {
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

func (c *Client) ConnectWebSocket() error {
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

func (c *Client) sendRequest(req *http.Request, responseBody interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	// TODO apikeyを入れる
	//req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	if res.StatusCode < http.StatusOK || res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unknown error, status code: %d", res.StatusCode)
	}

	if err = json.NewDecoder(res.Body).Decode(&responseBody); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
