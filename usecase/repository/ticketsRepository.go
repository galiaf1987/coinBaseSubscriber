package repository

import "github.com/galiaf1987/coinBaseSubscriber/usecase"

type TicketsRepository struct {
	BaseRepository
}

func (TicketsRepository) Save(ticks usecase.Ticks) bool {
	return true
}
