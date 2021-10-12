package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	coinbaseUrl = "wss://ws-feed.exchange.coinbase.com"
	origin      = "http://localhost/"
)

func SubscribeToCoinbaseMatches(c chan []byte, pairs []string) {
	ws, err := websocket.Dial(coinbaseUrl, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	subMsg := struct {
		Type       string   `json:"type"`
		ProductIDs []string `json:"product_ids"`
		Channels   []string `json:"channels"`
	}{
		Type:       "subscribe",
		ProductIDs: pairs,
		Channels:   []string{"matches"},
	}
	b, err := json.Marshal(subMsg)
	if _, err := ws.Write(b); err != nil {
		log.Fatal(fmt.Errorf("fail to write subscribe message: %v", err))
	}

	var isSubscribe bool
	for {
		var msg = make([]byte, 512)
		n, err := ws.Read(msg)
		if err != nil {
			log.Fatal(fmt.Errorf("fail to read incoming message: %v", err))
		}
		if !isSubscribe && strings.Contains(string(msg), "subscription") {
			isSubscribe = true
			continue
		}
		c <- msg[:n]
	}
}
