package listener

import (
	"fmt"
	"github.com/galiaf1987/coinBaseSubscriber/app/di"
	"github.com/galiaf1987/coinBaseSubscriber/usecase"
)

type Listener struct {
	provider usecase.RateProvider
}

func NewCoinBase(di di.DI) *Listener {
	return &Listener{
		provider: di.RateProvider,
	}
}

func (l *Listener) Listen() {
	tools := []string{"ETH-BTC", "BTC-USD", "BTC-EUR"}
	ch := make(chan usecase.Ticks)

	for _, tool := range tools {
		go l.subscribe(ch, tool)
	}

	for {
		fmt.Println(<-ch)
	}
}

func (l *Listener) subscribe(ch chan usecase.Ticks, tool string) {
	l.provider.GetRate(ch, tool)
}
