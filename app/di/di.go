package di

import (
	"github.com/galiaf1987/coinBaseSubscriber/environment"
	"github.com/galiaf1987/coinBaseSubscriber/usecase"
	"github.com/galiaf1987/coinBaseSubscriber/usecase/provider"
	"github.com/galiaf1987/coinBaseSubscriber/usecase/repository"
)

type DI struct {
	Config            environment.Config
	TicketsRepository usecase.TicksRepository
	RateProvider      usecase.RateProvider
}

func NewDI(cfg environment.Config) (di DI) {
	di = DI{
		Config:       cfg,
		RateProvider: provider.CoinBase{},
	}
	setupRepositoriesForDi(&di)

	return
}

func setupRepositoriesForDi(di *DI) {
	di.TicketsRepository = repository.TicketsRepository{}
}
