package world

import (
	"container/list"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

type CellType int

var StartingEnergy = 15

const (
	Empty CellType = iota
	Food
	AgentEn
	Obstacle
)

const (
	EnergyPerTick            = 1
	EnergyFoodGain           = 10
	EnergyReproduceThreshold = 12
	EnergyReproduceCost      = 5
)

type Cell struct {
	Type CellType
}

type World struct {
	Width      int
	Height     int
	AmountFood int
	TickCount  int
	DeathCount int
	BornCount  int
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
	DeathCount int
	BornCount  int
	AgentCount int
	PathPoints map[[2]int]bool
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

	world.BornCount = starterAgent

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
		prevX, prevY := a.X, a.Y
		nextX, nextY, found := w.FindTheClosestFood(a.X, a.Y, a)

		if found {
			a.ReduceEnergy()
		} else {
			nextX, nextY = a.MoveAiminglessly(w)
		}
		a.SetAgentPosition(nextX, nextY)
		if len(a.Path) > 0 {
			a.Path = a.Path[1:]
		}

		a.Eat(w)

		a.Die(w, 500*time.Millisecond)

		// Reflect Into World Map
		w.Grid[prevY][prevX].Type = Empty
		w.Grid[nextY][nextX].Type = AgentEn

		// newAgent := a.Reproduction(a.ID, w)
		// if newAgent != nil {
		// 	w.BornCount++
		// 	w.Agents = append(w.Agents, newAgent)
		// }
	}

	growth := rand.IntN(1000)
	if growth < 25 {
		w.spawnFood()
	}
}

func (w *World) RemoveAgent(target *Agent, duration time.Duration) {
	target.IsDie = true
	w.DeathCount++

	id := target.ID
	for i, a := range w.Agents {
		if a.ID == id {
			w.Agents = append(w.Agents[:i], w.Agents[i+1:]...)
			break
		}
	}

	go func(x, y int) {
		time.Sleep(duration)

		w.mu.Lock()
		defer w.mu.Unlock()
		w.Grid[y][x].Type = Empty
	}(target.X, target.Y)
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
		worldCopy.AgentCount++
	}

	avgEnergy := 0.0
	if len(w.Agents) > 0 {
		avgEnergy = float64(totalEnergy) / float64(len(w.Agents))
	}

	worldCopy.AvgEnergy = avgEnergy
	worldCopy.Tick = w.TickCount
	worldCopy.AmountFood = w.AmountFood
	worldCopy.BornCount = w.BornCount
	worldCopy.DeathCount = w.DeathCount

	worldCopy.PathPoints = make(map[[2]int]bool)

	for _, a := range worldCopy.Agents {
		for _, p := range a.Path {
			worldCopy.PathPoints[[2]int{p.x, p.y}] = true
		}
		fmt.Printf("Agent %d, Agent Location X: %d, Y:  %d Path %v\n", a.ID, a.X, a.Y, a.Path)
	}

	return worldCopy
}

type Chord struct {
	x, y int
}

func (w *World) FindTheClosestFood(currentX, currentY int, a *Agent) (int, int, bool) {
	visited := make([][]bool, w.Height)
	for i := range visited {
		visited[i] = make([]bool, w.Width)
	}

	q := list.New()
	q.PushBack(Chord{x: currentX, y: currentY})
	visited[currentY][currentX] = true

	parent := map[[2]int][2]int{}
	var target *Chord

	for q.Len() > 0 {
		cur := q.Remove(q.Front()).(Chord)

		if (cur.x != currentX || cur.y != currentY) && w.Grid[cur.y][cur.x].Type == Food {
			target = &cur
			break
		}

		for _, d := range DIRS {
			nx, ny := cur.x+d[0], cur.y+d[1]

			if nx < 0 || nx >= w.Width || ny < 0 || ny >= w.Height {
				continue
			}

			if visited[ny][nx] == true {
				continue
			}

			if w.Grid[ny][nx].Type == Obstacle {
				continue
			}

			visited[ny][nx] = true
			parent[[2]int{nx, ny}] = [2]int{cur.x, cur.y}
			q.PushBack(Chord{x: nx, y: ny})
		}
	}

	if target == nil {
		return currentX, currentY, false
	}

	px, py := target.x, target.y
	for !(px == currentX && py == currentY) {
		a.Path = append(a.Path, Chord{px, py})
		pxpy := parent[[2]int{px, py}]
		px, py = pxpy[0], pxpy[1]
	}

	// reverse path so it's from agent → food
	for i := 0; i < len(a.Path)/2; i++ {
		a.Path[i], a.Path[len(a.Path)-1-i] = a.Path[len(a.Path)-1-i], a.Path[i]
	}

	return a.Path[0].x, a.Path[0].y, true
}
