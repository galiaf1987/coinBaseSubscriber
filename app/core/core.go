package core

import (
	coinBaseListener "github.com/galiaf1987/coinBaseSubscriber/app/listener"
)

type Core struct {
	listener *coinBaseListener.Listener
}

func NewCore(
	listener *coinBaseListener.Listener,
) Core {
	return Core{
		listener: listener,
	}
}

func (c *Core) Run() {
	c.listener.Listen()
}
