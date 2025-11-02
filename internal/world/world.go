package world

import (
	"math/rand/v2"
	"sync"
	"time"
)

type CellType int

var StartingEnergy = 10

const (
	Empty CellType = iota
	Food
	AgentEn
	Obstacle
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
	Grid       [][]Cell
	Agents     []Agent
	Tick       int
	AvgEnergy  float64
	AmountFood int
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

	randomTotalObstacles := rand.IntN(width*height) - 300
	minTotalObstacles := randomTotalObstacles

	minTotalObstacles = max(50, minTotalObstacles)

	for range minTotalObstacles {
		randX := rand.IntN(width - 1)
		randY := rand.IntN(height - 1)

		world.Grid[randY][randX].Type = Obstacle
	}

	// Spawn Minim Food
	for range width * height / 5 {
		world.spawnFood()
	}

	freeCells := make([][2]int, 0)

	for y := range height {
		for x := range width {
			location := world.Grid[y][x].Type

			if location != Obstacle && location != Food {
				freeCells = append(freeCells, [2]int{x, y})
			}
		}
	}

	rand.Shuffle(len(freeCells), func(i, j int) {
		freeCells[i], freeCells[j] = freeCells[j], freeCells[i]
	})

	for i := range starterAgent {
		location := freeCells[i]
		x, y := location[0], location[1]
		agent := NewAgent(i, x, y, StartingEnergy)

		world.Agents = append(world.Agents, agent)
		world.Grid[y][x].Type = AgentEn
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

		newAgent := a.Reproduction(a.ID, w)
		if newAgent != nil {
			w.Agents = append(w.Agents, newAgent)
		}

		a.Die(w, 600*time.Millisecond)
	}

	if w.TickCount%10 == 0 {
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

	totalEnergy := 0
	worldCopy.Agents = make([]Agent, len(w.Agents))
	for i, a := range w.Agents {
		worldCopy.Agents[i] = *a // copy by value, not by pointer
		totalEnergy += a.Energy
	}

	avgEnergy := 0.0
	if len(w.Agents) > 0 {
		avgEnergy = float64(totalEnergy) / float64(len(w.Agents))
	}

	worldCopy.AvgEnergy = avgEnergy
	worldCopy.Tick = w.TickCount
	worldCopy.AmountFood = w.AmountFood

	return worldCopy
}
