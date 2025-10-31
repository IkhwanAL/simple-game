package world

import (
	"math/rand/v2"
	"sync"
	"time"
)

type CellType int

var StartingEnergy = 50

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

func NewWorld(width, height, starterAgent int) *World {
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

	for x := range starterAgent {
		randX := rand.IntN(width)
		ranxY := rand.IntN(height)

		agent := NewAgent(x, randX, ranxY, StartingEnergy)

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
	w.mu.Lock()

	defer w.mu.Unlock()

	w.TickCount++

	for _, a := range w.Agents {
		a.Move(w)
		a.Eat(w)

		newAgent := a.Reproduction(a.ID, w.Width-1, w.Height-1)
		if newAgent != nil {
			w.Agents = append(w.Agents, newAgent)
		}

		a.Die(w, 600*time.Millisecond)
	}

	if rand.Float64() < 0.1 {
		w.spawnFood()
	}

}

func (w *World) RemoveAgent(target *Agent, duration time.Duration) {
	target.IsDie = true

	go func(agentId int) {
		time.Sleep(duration)

		w.mu.Lock()
		defer w.mu.Unlock()

		w.RemoveAgentNow(agentId)

	}(target.ID)
}

func (w *World) RemoveAgentNow(agentId int) {
	for i, a := range w.Agents {
		if a.ID == agentId {
			w.Grid[a.Y][a.X].Type = Empty
			w.Agents = append(w.Agents[:i], w.Agents[i+1:]...)
		}
	}
}

func (w *World) Snapshot() WorldSnapshot {

	w.mu.RLock()
	defer w.mu.RUnlock()

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

	return worldCopy
}
