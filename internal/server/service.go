package server

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"sync/atomic"

	"time"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

type Service struct {
	world    *world.World
	stopChan chan struct{}
	ticker   *time.Ticker
	paused   atomic.Bool

	speedChan chan time.Duration
	Interval  time.Duration
}

func NewService(w *world.World) *Service {
	return &Service{
		speedChan: make(chan time.Duration),
		stopChan:  make(chan struct{}), // MUST initialize
		world:     w,
	}
}

func (s *Service) Snapshot() world.WorldSnapshot {
	return s.world.CaptureSnapshot()
}

func (s *Service) SpawnAgent() {
	s.world.SpawnAgent()
}

func (s *Service) SpawnFood() {
	s.world.SpawnFood()
}

func (s *Service) StartTick(interval time.Duration, hub *WebSocketHub) {
	if s.ticker != nil {
		s.ticker.Stop()
		for len(s.ticker.C) > 0 {
			<-s.ticker.C // drain any left over time
		}
	}

	s.Interval = interval
	s.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				if s.paused.Load() {
					continue
				}
				s.world.Tick()
				snapshot := s.Snapshot()

				msg, err := json.Marshal(snapshot)
				if err != nil {
					log.Fatal(err)
				}
				hub.Broadcast(context.Background(), msg)
			case newInterval := <-s.speedChan:
				s.ticker.Stop()
				for len(s.ticker.C) > 0 {
					<-s.ticker.C
				}

				s.Interval = newInterval
				s.ticker = time.NewTicker(newInterval)

				time.Sleep(interval)
			case <-s.stopChan:
				s.ticker.Stop()
				log.Println("The World Cease To Exists")
				return
			}
		}
	}()
}

func (s *Service) ChangeSpeed(tick time.Duration) {
	log.Printf("Change Speed To %s\n", tick)
	s.speedChan <- tick
}

func (s *Service) TogglePause() {
	log.Println("Calling Pause()")
	s.paused.Store(!s.paused.Load())
}

var stopOne sync.Once

func (s *Service) Stop() {
	stopOne.Do(func() {
		log.Println("Calling Stop()")
		close(s.stopChan)
	})

}
