# golang-limit-orderbook
A simple program written in Golang to implement Level 2 orderbook and matching algorithm

It creates 5 standalone orderbook instances and fill up 1 million order per order type per ticker via goroutines, totaling 10 million limit order objects. Afterwards, the program attempts to run matching executions simultaneously across tickers.

It runs about 45 seconds to finish both processes. 25 seconds dedicated to test limit order data filling and 20 seconds for matching executions.

Roughly 400K order insertions/sec and 500K order fulfillments/sec on M1 Pro CPU. Pretty good!