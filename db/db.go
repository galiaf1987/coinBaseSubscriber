package db

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"log"
	"net/url"
	"time"
)

type Options struct {
	Read            bool
	Write           bool
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
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=%s&charset=utf8mb4&collation=utf8mb4_bin",
		options.User,
		options.Password,
		options.Host,
		options.Port,
		options.Database,
		url.QueryEscape(options.Loc),
	))

	if err != nil {
		log.Fatalln("Не удалось инициализировать соединение к базе данных")
	}

	if options.MaxConn > 0 {
		connection.DB().SetMaxOpenConns(options.MaxConn)
	}

	if options.MaxIdleConn > 0 {
		connection.DB().SetMaxIdleConns(options.MaxIdleConn)
	}

	if options.ConnLifetimeSec > 0 {
		connection.DB().SetConnMaxLifetime(time.Duration(options.ConnLifetimeSec) * time.Second)
	}

	return connection
}
