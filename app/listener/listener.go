package listener

import (
	"github.com/galiaf1987/coinBaseSubscriber/app/di"
	"github.com/galiaf1987/coinBaseSubscriber/usecase"
	"time"
)

type Listener struct {
	provider        usecase.RateProvider
	ticksRepository usecase.TicksRepository
}

func NewCoinBase(di di.DI) *Listener {
	return &Listener{
		provider:        di.RateProvider,
		ticksRepository: di.TicketsRepository,
	}
}

var prevTicks map[string]usecase.Ticks

func (l Listener) Listen() {
	tools := []string{"ETH-BTC", "BTC-USD", "BTC-EUR"}
	channels := map[string]chan usecase.Ticks{
		"ETH-BTC": make(chan usecase.Ticks),
		"BTC-USD": make(chan usecase.Ticks),
		"BTC-EUR": make(chan usecase.Ticks),
	}

	for _, tool := range tools {
		go l.subscribe(channels[tool], tool)
	}

	entitiesForSave := make(map[string]usecase.Ticks, 3)

	for {
		select {
		case entitiesForSave["ETH-BTC"] = <-channels["ETH-BTC"]:
		case entitiesForSave["BTC-USD"] = <-channels["BTC-USD"]:
		case entitiesForSave["BTC-EUR"] = <-channels["BTC-EUR"]:
		}

		time.Sleep(time.Second)

		if len(entitiesForSave) != 0 {
			l.save(entitiesForSave)
		}
	}
}

func (l Listener) subscribe(ch chan usecase.Ticks, tool string) {
	l.provider.GetRate(ch, tool)
}

func (l Listener) save(data map[string]usecase.Ticks) bool {
	if prevTicks == nil {
		prevTicks = map[string]usecase.Ticks{}
	}

	var filteredData []usecase.Ticks
	for key, value := range data {
		if prevVal, ok := prevTicks[key]; ok && prevVal.Timestamp == value.Timestamp {
			continue
		}

		filteredData = append(filteredData, value)

		prevTicks[key] = value
	}

	if len(filteredData) > 0 {
		return l.ticksRepository.SaveMany(filteredData)
	}

	return false
}
