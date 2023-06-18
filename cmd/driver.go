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
	N            = 1_000_000
	BuyMaxPrice  = 300
	BuyMinPrice  = 150
	SellMaxPrice = 450
	SellMinPrice = 250
)

func main() {
	orderbooks := map[string]*entities.OrderBook{
		"aapl": entities.NewOrderBook("aapl"),
		"hood": entities.NewOrderBook("hood"),
		"spy":  entities.NewOrderBook("spy"),
		"shop": entities.NewOrderBook("shop"),
		"qqq":  entities.NewOrderBook("qqq"),
	}

	wg := &sync.WaitGroup{}

	fmt.Println("Populating limit orders")
	orderPopulateStart := time.Now()

	// Populate bid orders
	for ticker, ob := range orderbooks {
		wg.Add(1)
		go addTestLimitOrder(wg, ob, ticker, true)
	}

	// Populate ask orders
	for ticker, ob := range orderbooks {
		wg.Add(1)
		go addTestLimitOrder(wg, ob, ticker, false)
	}

	wg.Wait()
	fmt.Printf("Time to populate limit orders: %d seconds", (time.Now().UnixMilli()-orderPopulateStart.UnixMilli())/int64(1000))

	for ticker, ob := range orderbooks {
		fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())
	}

	orderMatchingStart := time.Now()
	fmt.Println("Order matching begins...")
	for ticker, ob := range orderbooks {
		wg.Add(1)
		go executeMatching(ob, ticker, wg)
	}

	wg.Wait()
	fmt.Printf("Time to finish order matching: %d seconds", (time.Now().UnixMilli()-orderMatchingStart.UnixMilli())/int64(1000))
}

func executeMatching(ob *entities.OrderBook, ticker string, wg *sync.WaitGroup) {
	ob.Match()
	fmt.Printf("ticker: %s bids: %d asks: %d \n", ticker, ob.Bids.Len(), ob.Asks.Len())

	wg.Done()
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

		//fmt.Printf("Created limit order - ticker: %s price: %s quantity: %d isBid: %t \n", ticker, priceString, quantity, isBid)
	}

	wg.Done()
}
