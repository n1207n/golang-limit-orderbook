package entities

import (
	utils "github.com/n1207n/golang-limit-orderbook/utils"
)

type OrderBook struct {
	bids utils.PriorityQueue
	asks utils.PriorityQueue
}
