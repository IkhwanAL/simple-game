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
	crtl     *world.WorldController
	stopChan chan struct{}
	ticker   *time.Ticker
	paused   atomic.Bool

	speedChan chan time.Duration
	Interval  time.Duration
}

func NewService(ctrl *world.WorldController) *Service {
	return &Service{
		speedChan: make(chan time.Duration),
		stopChan:  make(chan struct{}), // MUST initialize
		crtl:      ctrl,
	}
}

func (s *Service) Snapshot() world.WorldSnapshot {
	reply := make(chan world.WorldSnapshot)
	s.crtl.CmdChan <- world.CmdSnapshot{Reply: reply}
	return <-reply
}

func (s *Service) SpawnAgent() {
	s.crtl.CmdChan <- world.CmdSpawnAgent{}
}

func (s *Service) SpawnFood() {
	s.crtl.CmdChan <- world.CmdSpawnFood{}
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
				s.crtl.CmdChan <- world.CmdTick{}

				snapshot := s.Snapshot()
				msg, err := json.Marshal(snapshot)
				if err != nil {
					log.Fatal(err)
				}
				hub.Broadcast(context.Background(), msg)
			case newInterval := <-s.speedChan:
				s.Interval = newInterval

				s.ticker.Reset(newInterval)

			case <-s.stopChan:
				s.ticker.Stop()
				s.crtl.Stop()
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
		s.crtl.Stop()
		log.Println("Calling Stop()")
		close(s.stopChan)
	})
}
