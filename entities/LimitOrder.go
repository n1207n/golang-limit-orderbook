package entities

import (
	uuid "github.com/google/uuid"
)

type LimitOrder struct {
	ID       uuid.UUID
	Ticker   string
	Price    int
	Quantity int
}
