package utils

import (
	"fmt"
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
}

func (lo LimitOrder) String() string {
	return fmt.Sprintf("[%s - %s - %d - Is Buy: %t - %d - %s]\n", lo.Ticker, lo.Price.String(), lo.Quantity, lo.IsBid, lo.Timestamp.UnixMilli(), lo.ID.String())
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
}

func (pq *OrderPriorityQueue) Push(x any) {
	limitOrder := x.(*LimitOrder)
	*pq = append(*pq, limitOrder)
}

func (pq *OrderPriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	limitOrder := old[n-1]
	*pq = old[0 : n-1]
	return limitOrder
}

func (pq *OrderPriorityQueue) Peek() any {
	old := *pq
	return old[0]
}
