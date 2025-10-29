package world

import "sync"

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
