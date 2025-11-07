package world

import (
	"container/list"
	"math/rand/v2"
	"sync"
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
	Width       int
	Height      int
	AmountFood  int
	TickCount   int
	DeathCount  int
	BornCount   int
	Grid        [][]Cell
	Agents      []*Agent
	PendingDead []*Agent
	Mu          sync.RWMutex
	DebugMode   bool
}

func NewWorld(width, height, starterAgent int, isDebugOn bool) *World {
	world := &World{
		Height:    height,
		Width:     width,
		DebugMode: isDebugOn,
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
		world.SpawnFood()
	}
	// for range 4 {
	// 	world.SpawnFood()
	// }

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
		agent := NewAgent(x, y, StartingEnergy)

		world.Agents = append(world.Agents, agent)
		world.Grid[y][x].Type = AgentEn
	}

	world.BornCount = starterAgent

	return world
}

func (w *World) AddAgent(a *Agent) {
	w.Agents = append(w.Agents, a)
}

func (w *World) SpawnFood() {
	x, y := rand.IntN(w.Height), rand.IntN(w.Width)

	if w.Grid[y][x].Type == Empty {
		w.Grid[y][x].Type = Food
		w.AmountFood += 1
	}
}

func (w *World) Tick() {
	w.Mu.Lock()

	defer w.Mu.Unlock()

	w.TickCount++

	if len(w.PendingDead) > 0 {
		for _, deadAgent := range w.PendingDead {
			w.Grid[deadAgent.Y][deadAgent.X].Type = Empty
			w.RemoveAgent(deadAgent)
		}

		w.PendingDead = nil
	}

	for _, a := range w.Agents {
		if a.IsDie {
			continue
		}

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

		a.Die(w)

		// Reflect Into World Map
		w.Grid[prevY][prevX].Type = Empty
		w.Grid[nextY][nextX].Type = AgentEn

		newAgent := a.Reproduction(w)
		if newAgent != nil {
			w.BornCount++
			w.Agents = append(w.Agents, newAgent)
		}
	}

	growth := rand.IntN(1000)
	if growth < 25 {
		w.SpawnFood()
	}
}

func (w *World) RemoveAgent(target *Agent) {
	id := target.ID
	for i, a := range w.Agents {
		if a.ID == id {
			w.Agents = append(w.Agents[:i], w.Agents[i+1:]...)
			break
		}
	}

}

type WorldSnapshot struct {
	Tick       int             `json:"tick"`
	Width      int             `json:"width"`
	Height     int             `json:"height"`
	Food       [][2]int        `json:"foods"`
	Agents     []AgentSnapshot `json:"agents"`
	Obstacle   [][2]int        `json:"obstacles"`
	BornCount  int             `json:"bornCount"`
	DeathCount int             `json:"deathCount"`
}

type AgentSnapshot struct {
	ID     int  `json:"id"`
	X      int  `json:"x"`
	Y      int  `json:"y"`
	IsDead bool `json:"isDead"`
}

func (w *World) Snapshot() WorldSnapshot {
	w.Mu.RLock()
	defer w.Mu.RUnlock()

	var aCopy WorldSnapshot

	aCopy.Tick = w.TickCount
	aCopy.Width = w.Width
	aCopy.Height = w.Height
	aCopy.BornCount = w.BornCount
	aCopy.DeathCount = w.DeathCount

	var foods [][2]int
	var obstacles [][2]int

	for y, row := range w.Grid {
		for x, column := range row {
			switch column.Type {
			case Food:
				foods = append(foods, [2]int{x, y})
			case Obstacle:
				obstacles = append(obstacles, [2]int{x, y})
			}
		}
	}

	aCopy.Food = foods

	if aCopy.Food == nil {
		aCopy.Food = [][2]int{}
	}

	aCopy.Obstacle = obstacles

	var agents []AgentSnapshot

	for _, a := range w.Agents {
		agents = append(agents, AgentSnapshot{
			ID:     a.ID,
			X:      a.X,
			Y:      a.Y,
			IsDead: a.IsDie,
		})
	}

	aCopy.Agents = agents
	if aCopy.Agents == nil {
		aCopy.Agents = []AgentSnapshot{}
	}

	return aCopy
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

			if w.Grid[ny][nx].Type == Obstacle || w.Grid[ny][nx].Type == AgentEn {
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
