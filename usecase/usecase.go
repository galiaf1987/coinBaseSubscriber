package usecase

import "time"

type Ticks struct {
	Timestamp time.Time
	Symbol    string
	BestBid   float64
	BestAsk   float64
}

type TicksRepository interface {
	Save(ticks Ticks) bool
}

type RateProvider interface {
	GetRate(ch chan Ticks)
}
