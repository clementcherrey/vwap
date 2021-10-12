package main

import (
	"encoding/json"
	"flag"
	"log"
	"regexp"
	"strings"
)

func main() {
	pairList := flag.String("pairs", "BTC-USD,ETH-USD,ETH-BTC", "trading pairs to monitor")
	windowSize := flag.Int("window", 200, "VWAP window size")
	flag.Parse()

	// validate the input (list of the pairs)
	r, _ := regexp.Compile(`^([A-Z]{3}\-[A-Z]{3},)*([A-Z]{3}\-[A-Z]{3})$`)
	if !r.MatchString(*pairList) {
		log.Fatal("invalid pair list")
	}
	pairNames := strings.Split(*pairList, ",")

	// create an aggregator for each pair and a channel to send incoming sizedPrices to it
	pairs := make(map[string]chan *match)
	for _, name := range pairNames {
		aggregator := NewPairAggregator(name, *windowSize)
		incomingMatches := make(chan *match)
		go aggregator.ListenForNewMatch(incomingMatches)
		pairs[name] = incomingMatches
	}

	// subscribe to "matches" for the list of pairs
	incomingMessages := make(chan []byte)
	go SubscribeToCoinbaseMatches(incomingMessages, pairNames)

	// send incoming match to the correct aggregator
	// by using the channel that was associated to it
	for msg := range incomingMessages {
		m := &match{}
		if err := json.Unmarshal(msg, m); err != nil {
			log.Fatalf("fail to Unmarshal %v", err)
		}

		c, ok := pairs[m.ProductID]
		if !ok {
			log.Println(m)
			log.Fatalf("no pair found for %s", m.ProductID)
		}
		c <- m
	}
}
