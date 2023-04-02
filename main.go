package main

import (
	"flag"
	"github.com/galiaf1987/coinBaseSubscriber/app/core"
	"github.com/galiaf1987/coinBaseSubscriber/app/di"
	"github.com/galiaf1987/coinBaseSubscriber/app/listener"
	"github.com/galiaf1987/coinBaseSubscriber/environment"
	"log"
)

func main() {
	configPath := flag.String(
		"config",
		"config/development/config.toml",
		"Путь до файла конфига toml")
	flag.Parse()

	cfg := environment.NewConfig(*configPath)
	if err := cfg.Init(); err != nil {
		log.Fatalln("unable to parse config", err)
	}

	appCore := initCore(cfg)

	appCore.Run()
}

func initCore(cfg environment.Config) *core.Core {
	newDI := di.NewDI(cfg)
	appCore := core.NewCore(listener.NewCoinBase(newDI))

	return &appCore
}
