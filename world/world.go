package world

import (
	// "fmt"
	"math/rand/v2"
	"sync"
)

const (
	Width  = 20
	Height = 20
)

type CellType int

const (
	Empty CellType = iota
	Food
	AgentEn
)

type Cell struct {
	Type CellType
}

type World struct {
	Grid   [Height][Width]Cell
	mu     sync.RWMutex
	Agents []*Agent
}

type WorldViewData struct {
	Grid   [Height][Width]Cell
	Agents []*Agent
}

func NewWorld() *World {
	world := &World{}

	for range 30 {
		world.spawnFood()
	}

	return world
}

func (w *World) AddAgent(a *Agent) {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.Agents = append(w.Agents, a)
}

func (w *World) spawnFood() {
	x, y := rand.IntN(Height), rand.IntN(Width)

	if w.Grid[y][x].Type == Empty {
		w.Grid[y][x].Type = Food
	}
}

func (w *World) Tick() {
	w.mu.Lock()
	defer w.mu.Unlock()

	for _, a := range w.Agents {
		a.Act(w)
	}

	if rand.Float64() < 0.2 {
		w.spawnFood()
	}
}

func (w *World) Snapshot() WorldViewData {
	w.mu.RLock()
	defer w.mu.RUnlock()

	var copy WorldViewData

	copy.Grid = w.Grid
	copy.Agents = append([]*Agent(nil), w.Agents...)

	return copy
}
