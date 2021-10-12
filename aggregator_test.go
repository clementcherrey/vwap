package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPairAggregator_add(t *testing.T) {
	pa := pairAggregator{
		sizedPrices:              make([]*sizedPrice, 101),
		totalSize:                289.1,
		totalVolumeWeightedPrice: 4321.5,
		windowSize:               200,
	}
	pa.add(&sizedPrice{
		size:  7.2,
		price: 20.9,
	})
	assert.Equal(t, 296.3, pa.totalSize)
	assert.Equal(t, 4321.5+(7.2*20.9), pa.totalVolumeWeightedPrice)
	assert.Equal(t, 200, pa.windowSize)
}

func TestToSizedPrice(t *testing.T) {
	_, err := toSizedPrice(&match{
		Size:  "0",
		Price: "123",
	})
	assert.Error(t, err)
	assert.EqualError(t, err, "invalid price (123) or size (0)")

	_, err = toSizedPrice(&match{
		Size:  "1.111111111111111",
		Price: "1e6",
	})
	assert.NoError(t, err)
}
