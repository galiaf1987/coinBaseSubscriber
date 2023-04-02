package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
)

type Options struct {
	Host            string
	Port            int
	User            string
	Password        string
	Database        string
	MaxConn         int
	MaxIdleConn     int
	ConnLifetimeSec int
	Loc             string
}

func GetConnection(options Options) *gorm.DB {
	connection, err := gorm.Open("mysql", fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&&charset=utf8mb4&collation=utf8mb4_bin",
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
	))

	if err != nil {
		log.Fatalln("Не удалось инициализировать соединение к базе данных")
	}

	return connection
}
