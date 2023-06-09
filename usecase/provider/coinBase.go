package provider

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/galiaf1987/coinBaseSubscriber/usecase"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

type Options struct {
	Url string
}

type CoinBase struct {
	Options Options
}

type SubscribeRequest struct {
	Type       string   `json:"type"`
	ProductIds []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type TickerResponse struct {
	Type        string    `json:"type"`
	Sequence    int64     `json:"sequence"`
	ProductId   string    `json:"product_id"`
	Price       string    `json:"price"`
	Open24H     string    `json:"open_24h"`
	Volume24H   string    `json:"volume_24h"`
	Low24H      string    `json:"low_24h"`
	High24H     string    `json:"high_24h"`
	Volume30D   string    `json:"volume_30d"`
	BestBid     string    `json:"best_bid"`
	BestBidSize string    `json:"best_bid_size"`
	BestAsk     string    `json:"best_ask"`
	BestAskSize string    `json:"best_ask_size"`
	Side        string    `json:"side"`
	Time        time.Time `json:"time"`
	TradeId     int       `json:"trade_id"`
	LastSize    string    `json:"last_size"`
}

func (cb CoinBase) GetRate(ch chan usecase.Ticks, tool string) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: cb.Options.Url}

	c := connect(u)
	defer c.Close()

	subscribe(c, tool)

	go func() {
		defer func() {
			fmt.Println("reconnecting")
			if e := recover(); e != nil {
				c = connect(u)
			}
		}()
		for {
			readMessage(c, ch)
		}
	}()

	for {
		select {
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			return
		}
	}
}

func connect(u url.URL) *websocket.Conn {
	log.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	return c
}

func readMessage(c *websocket.Conn, ch chan usecase.Ticks) bool {
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Printf("Read error: %s", err)
		return false
	}

	response := TickerResponse{}
	err = json.Unmarshal(message, &response)

	if err != nil {
		log.Printf("incorrect response: %s", message)
		return false
	}

	if response.Type != "ticker" {
		return false
	}

	if err != nil {
		log.Printf("Convertation error: %s", err)
		return false
	}

	domainEntity := usecase.Ticks{
		Timestamp: response.Time,
		Symbol:    response.ProductId,
		BestBid:   response.BestBid,
		BestAsk:   response.BestAsk,
	}

	ch <- domainEntity
	return true
}

func subscribe(c *websocket.Conn, tool string) {
	request := SubscribeRequest{
		Type:       "subscribe",
		ProductIds: []string{tool},
		Channels:   []string{"ticker"},
	}

	err := c.WriteJSON(request)
	if err != nil {
		panic("Invalid Request")
	}
}
