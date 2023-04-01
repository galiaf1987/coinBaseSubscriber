package repository

import (
	"fmt"
	"github.com/galiaf1987/coinBaseSubscriber/usecase"
)

type TicketsRepository struct {
	BaseRepository
}

func (TicketsRepository) SaveMany(ticks []usecase.Ticks) bool {
	fmt.Println(ticks)
	return true
}
