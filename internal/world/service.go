package world

import (
	"sync"
	"time"
)

type Service struct {
	mu    sync.Mutex
	world *World
}

func NewService(w *World) *Service {
	return &Service{world: w}
}

func (s *Service) Tick() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.world.Tick()
}

func (s *Service) Snapshot() WorldSnapshot {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.world.Snapshot()
}

func (s *Service) StartTick(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			s.Tick()
			Logf("Tick Completed") // Will Be Deleted After Clarification
		}
	}()

	Logf("World is Starting")
}
