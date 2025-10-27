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
	mu     sync.Mutex
	Agents []*Agent
}

func NewWorld() *World {
	world := &World{}

	for range 30 {
		world.spawnFood()
	}

	return world
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
		// fmt.Printf("%s\t", "Action From Struct Agent Called")
		a.Act(w)
	}

	if rand.Float64() < 0.2 {
		w.spawnFood()
	}
}
