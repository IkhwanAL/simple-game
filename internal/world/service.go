package world

import (
	"log"
	"math/rand"
	"sync"
	"sync/atomic"

	"time"
)

type Service struct {
	world    *World
	stopChan chan struct{}
	ticker   *time.Ticker
	paused   atomic.Bool

	speedChan chan time.Duration
	Interval  time.Duration
}

func NewService(w *World) *Service {
	return &Service{
		speedChan: make(chan time.Duration),
		stopChan:  make(chan struct{}), // MUST initialize
		world:     w,
	}
}

func (s *Service) Snapshot() WorldSnapshot {
	return s.world.Snapshot()
}

func (s *Service) SpawnAgent() {
	s.world.mu.Lock()
	defer s.world.mu.Unlock()

	randomNum := rand.Intn(999)

	agent := NewAgent(randomNum, rand.Intn(s.world.Width-1), rand.Intn(s.world.Height-1), StartingEnergy)
	s.world.AddAgent(agent)
}

func (s *Service) SpawnFood() {
	s.world.mu.Lock()
	defer s.world.mu.Unlock()

	s.world.spawnFood()
}

func (s *Service) StartTick(interval time.Duration) {
	if s.ticker != nil {
		s.ticker.Stop()
		for len(s.ticker.C) > 0 {
			<-s.ticker.C // drain zombie ticks
		}
	}

	s.Interval = interval
	s.ticker = time.NewTicker(interval)

	go func() {
		for {
			select {
			case <-s.ticker.C:
				// fmt.Println("Ticker Called", time.Now().Format(time.RFC3339Nano), "svc=", s)
				// fmt.Println("Ticker Called", s.ticker, time.Now())
				if s.paused.Load() {
					continue
				}
				s.world.Tick()
			case newInterval := <-s.speedChan:
				log.Printf("Try To Change Speed %s", newInterval)
				s.ticker.Stop()
				s.Interval = newInterval
				s.ticker = time.NewTicker(newInterval)
				log.Printf("Change Tick Speed %s", newInterval)
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
