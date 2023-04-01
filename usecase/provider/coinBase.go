package provider

import (
	"encoding/json"
	"flag"
	"github.com/galiaf1987/coinBaseSubscriber/usecase"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

type Options struct {
	Url string
}

type CoinBase struct{}

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

var addr = flag.String("addr", "ws-feed-public.sandbox.exchange.coinbase.com", "http service address")

func (CoinBase) GetRate(ch chan usecase.Ticks, tool string) {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	u := url.URL{Scheme: "wss", Host: *addr}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	subscribe(c, tool)

	go func() {
		defer close(done)
		for {
			readMessage(c, ch)
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-interrupt:
			log.Println("interrupt")

			// Cleanly close the connection by sending a close message and then
			// waiting (with timeout) for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func readMessage(c *websocket.Conn, ch chan usecase.Ticks) bool {
	_, message, err := c.ReadMessage()
	if err != nil {
		log.Printf("Read error: %s", message)
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

	bestBid, err := strconv.ParseFloat(response.BestBid, 64)
	bestAsk, err := strconv.ParseFloat(response.BestAsk, 64)
	if err != nil {
		log.Printf("Convertation error: %s", err)
		return false
	}

	domainEntity := usecase.Ticks{
		Timestamp: response.Time,
		Symbol:    response.ProductId,
		BestBid:   bestBid,
		BestAsk:   bestAsk,
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
