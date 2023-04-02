package repository

import (
	"fmt"
	"github.com/galiaf1987/coinBaseSubscriber/usecase"
	"log"
)

type TicketsRepository struct {
	BaseRepository
}

type Ticks struct {
	Timestamp int64
	Symbol    string
	BestBid   string `gorm:"column:best_bid"`
	BestAsk   string `gorm:"column:best_ask"`
}

func (tr TicketsRepository) SaveMany(domainEntities []usecase.Ticks) bool {
	fmt.Println(domainEntities)

	for _, val := range domainEntities {
		ticks := Ticks{
			Timestamp: val.Timestamp.UnixMilli(),
			Symbol:    val.Symbol,
			BestBid:   val.BestBid,
			BestAsk:   val.BestAsk,
		}

		if err := tr.DBConnection.Create(&ticks).Error; err != nil {
			log.Fatal(err)

			return false
		}
	}

	tr.DBConnection.QueryExpr()
	return true
}
