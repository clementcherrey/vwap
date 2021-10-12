package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
)

// the minimal version of a match, as it is returned by the websocket
type match struct {
	ProductID string `json:"product_id"`
	Size      string `json:"size"`
	Price     string `json:"price"`
}

// a price with the associated size
type sizedPrice struct {
	size  float64
	price float64
}

func toSizedPrice(m *match) (*sizedPrice, error) {
	size, err := strconv.ParseFloat(m.Size, 64)
	if err != nil {
		return nil, errors.New("fail to parse match size")
	}
	price, err := strconv.ParseFloat(m.Price, 64)
	if err != nil {
		return nil, errors.New("fail to parse match price")
	}
	if price <= 0 || size <= 0 {
		return nil, fmt.Errorf("invalid price (%v) or size (%v)", price, size)
	}
	return &sizedPrice{
		size:  size,
		price: price,
	}, nil
}

type pairAggregator struct {
	pairName                            string
	sizedPrices                         []*sizedPrice
	totalSize, totalVolumeWeightedPrice float64
	windowSize                          int
}

func NewPairAggregator(pairName string, windowSize int) *pairAggregator {
	return &pairAggregator{
		pairName:   pairName,
		windowSize: windowSize,
	}
}

func (pa *pairAggregator) add(sp *sizedPrice) {
	pa.sizedPrices = append(pa.sizedPrices, sp)
	pa.totalSize += sp.size
	pa.totalVolumeWeightedPrice += sp.size * sp.price
}

func (pa *pairAggregator) removeOldest() {
	oldestMatch := pa.sizedPrices[0]
	pa.totalSize -= oldestMatch.size
	pa.totalVolumeWeightedPrice -= oldestMatch.size * oldestMatch.size
}

func (pa *pairAggregator) vwap() float64 {
	return pa.totalVolumeWeightedPrice / pa.totalSize
}

func (pa *pairAggregator) printVWAP() {
	log.Printf("%s VWAP %f\n", pa.pairName, pa.vwap())
}

func (pa *pairAggregator) update(m *match) {
	if len(pa.sizedPrices) == pa.windowSize {
		pa.removeOldest()
	}
	sp, err := toSizedPrice(m)
	if err != nil {
		panic(err)
	}
	pa.add(sp)
}

func (pa *pairAggregator) ListenForNewMatch(c chan *match) {
	for m := range c {
		pa.update(m)
		pa.printVWAP()
	}
}
