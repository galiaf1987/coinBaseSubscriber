package usecase

import "time"

type Ticks struct {
	Timestamp time.Time
	Symbol    string
	BestBid   string
	BestAsk   string
}

type TicksRepository interface {
	SaveMany(domainEntities []Ticks) bool
}

type RateProvider interface {
	GetRate(ch chan Ticks, tool string)
}
