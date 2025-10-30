package world

import (
	"log"

	"time"
)

type Service struct {
	world    *World
	stopChan chan struct{}
	ticker   time.Ticker
}

func NewService(w *World) *Service {
	return &Service{
		stopChan: make(chan struct{}), // MUST initialize
		world:    w,
	}
}

func (s *Service) Snapshot() WorldSnapshot {
	return s.world.Snapshot()
}

func (s *Service) StartTick(interval time.Duration) {
	s.ticker = *time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-s.ticker.C:
				log.Println("Tick")
				s.world.Tick()
			case <-s.stopChan:
				log.Println("The World Cease To Exists")
				return
			}
		}
	}()
}

func (s *Service) Stop() {
	log.Println("Calling Stop()")
	close(s.stopChan)
}
