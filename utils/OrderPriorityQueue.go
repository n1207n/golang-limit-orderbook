package utils

import (
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"time"
)

type LimitOrder struct {
	ID        uuid.UUID
	Ticker    string
	Price     decimal.Decimal
	Quantity  int
	IsBid     bool
	Timestamp time.Time
	Index     int // The index of the LimitOrder in the heap.
}

// A OrderPriorityQueue implements heap.Interface and holds LimitOrders.
type OrderPriorityQueue []*LimitOrder

func (pq OrderPriorityQueue) Len() int { return len(pq) }

func (pq OrderPriorityQueue) Less(i, j int) bool {
	// For bids, behave as max-heap
	// For asks, behave as min-heap
	if pq[i].IsBid {
		return pq[i].Price.GreaterThan(pq[j].Price)
	}

	return pq[i].Price.LessThan(pq[j].Price)
}

func (pq OrderPriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *OrderPriorityQueue) Push(x any) {
	n := len(*pq)
	limitOrder := x.(*LimitOrder)
	limitOrder.Index = n
	*pq = append(*pq, limitOrder)
}

func (pq *OrderPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	limitOrder := old[n-1]
	old[n-1] = nil        // avoid memory leak
	limitOrder.Index = -1 // for safety
	*pq = old[0 : n-1]
	return limitOrder
}
