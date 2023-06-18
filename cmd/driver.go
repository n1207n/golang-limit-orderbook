package main

import (
	"fmt"
	"github.com/n1207n/golang-limit-orderbook/entities"
	"math/rand"
	"strconv"
	"sync"
	"time"
)

const (
	N            = 100_000
	BuyMaxPrice  = 300
	BuyMinPrice  = 150
	SellMaxPrice = 450
	SellMinPrice = 250
)

func main() {
	orderbooks := map[string]*entities.OrderBook{
		"aapl": entities.NewOrderBook("aapl"),
		"hood": entities.NewOrderBook("hood"),
	}

	wg := &sync.WaitGroup{}

	fmt.Println("Populating limit orders")

	// Populate bid orders
	for ticker, ob := range orderbooks {
		wg.Add(1)
		addTestLimitOrder(wg, ob, ticker, true)
	}

	wg.Wait()

	// Populate ask orders
	for ticker, ob := range orderbooks {
		wg.Add(1)
		addTestLimitOrder(wg, ob, ticker, false)
	}

	wg.Wait()

	for ticker, ob := range orderbooks {
		fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())
	}

	fmt.Println("Order matching begins...")
	for ticker, ob := range orderbooks {
		ob.Match()
		fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())
	}
}

func addTestLimitOrder(wg *sync.WaitGroup, ob *entities.OrderBook, ticker string, isBid bool) {
	for i := 0; i < N; i++ {
		var randomPrice float64
		r := rand.New(rand.NewSource(time.Now().UnixNano()))

		if isBid {
			randomPrice = BuyMinPrice + r.Float64()*(BuyMaxPrice-BuyMinPrice)

		} else {
			randomPrice = SellMinPrice + r.Float64()*(SellMaxPrice-SellMinPrice)
		}

		priceString := strconv.FormatFloat(randomPrice, 'f', 2, 64)
		quantity := r.Intn(50000)

		ob.AddLimitOrder(
			ticker,
			priceString,
			quantity,
			isBid)

		fmt.Printf("Created limit order - ticker: %s price: %s quantity: %d isBid: %t \n", ticker, priceString, quantity, isBid)
	}

	wg.Done()
}
