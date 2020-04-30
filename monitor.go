package main

import (
	"go.mongodb.org/mongo-driver/event"
)

type poolMonitor struct {
	conns int
}

func newPoolMonitor() *poolMonitor {
	return &poolMonitor{}
}

func (p *poolMonitor) HandleEvent(evt *event.PoolEvent) {
	switch evt.Type {
	case event.ConnectionCreated:
		p.conns++
	case event.ConnectionClosed:
		p.conns--
	}
}
