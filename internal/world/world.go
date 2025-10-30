package world

import (
	"log"
	"math/rand/v2"
	"sync"
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
	Width      int
	Height     int
	AmountFood int
	TickCount  int
	Grid       [][]Cell
	Agents     []*Agent
	mu         sync.RWMutex
}

type WorldSnapshot struct {
	Grid   [][]Cell
	Agents []Agent
}

func NewWorld(width, height int) *World {
	world := &World{
		Height: height,
		Width:  width,
	}

	world.Grid = make([][]Cell, height)

	// Generate Grid World Map
	for x := range height {
		world.Grid[x] = make([]Cell, width)
	}

	// Spawn Minim Food
	for range 20 {
		world.spawnFood()
	}

	for x := range 5 {
		randX := rand.IntN(width)
		ranxY := rand.IntN(height)

		agent := NewAgent(x, randX, ranxY, 50)

		world.Agents = append(world.Agents, agent)
		world.Grid[ranxY][randX].Type = AgentEn
	}

	return world
}

func (w *World) AddAgent(a *Agent) {
	w.Agents = append(w.Agents, a)
}

func (w *World) spawnFood() {
	x, y := rand.IntN(w.Height), rand.IntN(w.Width)

	if w.Grid[y][x].Type == Empty {
		w.Grid[y][x].Type = Food
		w.AmountFood += 1
	}
}

func (w *World) Tick() {
	log.Println("Tick Waiting To Lock")
	w.mu.Lock()
	log.Println("Tick is Locking")

	defer w.mu.Unlock()

	w.TickCount++

	for _, a := range w.Agents {
		a.Move(w)
		a.Eat(w)
	}

	if rand.Float64() < 0.1 {
		w.spawnFood()
	}

	log.Println("Tick is Unlock")
}

func (w *World) Snapshot() WorldSnapshot {

	log.Println("Snapshot Waiting To RLock")
	w.mu.RLock()
	log.Println("Snapshot is RLocking")
	defer w.mu.RUnlock()

	log.Print("Snap")
	var worldCopy WorldSnapshot

	worldCopy.Grid = make([][]Cell, len(w.Grid))
	for i := range w.Grid {
		row := make([]Cell, len(w.Grid[i]))
		copy(row, w.Grid[i])
		worldCopy.Grid[i] = row
	}

	worldCopy.Agents = make([]Agent, len(w.Agents))
	for i, a := range w.Agents {
		worldCopy.Agents[i] = *a // copy by value, not by pointer
	}

	log.Println("Snapshot is RUnlock")
	return worldCopy
}
