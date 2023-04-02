package environment

import (
	"github.com/BurntSushi/toml"
	"github.com/galiaf1987/coinBaseSubscriber/db"
	"github.com/galiaf1987/coinBaseSubscriber/usecase/provider"
)

type Config struct {
	file               string
	Database           db.Options
	CoinBaseSubscriber provider.Options
}

func NewConfig(file string) Config {
	return Config{
		file: file,
	}
}

func (c *Config) Init() error {
	_, err := toml.DecodeFile(c.file, c)
	return err
}
