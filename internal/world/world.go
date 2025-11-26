package world

import (
	"math/rand/v2"
)

type CellType int

var StartingEnergy = 15

const (
	Empty CellType = iota
	Food
	AgentEn
	Obstacle
	BuffIncreaseFoV
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

type WorldCommand struct {
	Action string
	Data   any
	Reply  chan any
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
	DebugMode   bool
}

type InitWorld struct {
	Width         int
	Height        int
	StarterAgent  int
	IsDebugOn     bool // Remove All Obstacle
	totalInitBuff int
}

func NewWorld(init InitWorld) *World {
	height := init.Height
	width := init.Width

	world := &World{
		Height:    height,
		Width:     width,
		DebugMode: init.IsDebugOn,
	}

	world.Grid = make([][]Cell, height)

	// Generate Grid World Map
	for x := range height {
		world.Grid[x] = make([]Cell, width)
	}

	randomTotalObstacles := rand.IntN(width*height) - 300
	minTotalObstacles := randomTotalObstacles

	minTotalObstacles = max(50, minTotalObstacles)

	if !world.DebugMode {
		for range minTotalObstacles {
			randX := rand.IntN(width - 1)
			randY := rand.IntN(height - 1)

			world.Grid[randY][randX].Type = Obstacle
		}
	}

	for range width * height / 5 {
		world.SpawnFood()
	}

	for range init.totalInitBuff {
		world.SpawnBuff()
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

	for i := range init.StarterAgent {
		location := freeCells[i]
		x, y := location[0], location[1]
		agent := NewAgent(x, y, StartingEnergy)

		world.Agents = append(world.Agents, agent)
		world.Grid[y][x].Type = AgentEn
	}

	world.BornCount = init.StarterAgent

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

func (w *World) SpawnBuff() {
	x, y := rand.IntN(w.Height), rand.IntN(w.Width)

	if w.Grid[y][x].Type == Empty {
		w.Grid[y][x].Type = BuffIncreaseFoV
	}
}

func (w *World) Tick() {
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

		var nextX, nextY int

		a.TraitControl()

		act := a.ChooseAction()

		nextX, nextY = a.PerformAction(w, act)

		if nextX != prevX || nextY != prevY {
			a.ReduceEnergy()
			a.SetAgentPosition(nextX, nextY)
		}

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
	if growth < 250 {
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
	Tick         int             `json:"tick"`
	Width        int             `json:"width"`
	Height       int             `json:"height"`
	Food         [][2]int        `json:"foods"`
	AvgEnergy    float64         `json:"avgEnergy"`
	Agents       []AgentSnapshot `json:"agents"`
	Obstacle     [][2]int        `json:"obstacles"`
	BornCount    int             `json:"bornCount"`
	DeathCount   int             `json:"deathCount"`
	TickInterval int64           `json:"tickInterval"`
}

type AgentSnapshot struct {
	ID     int  `json:"id"`
	X      int  `json:"x"`
	Y      int  `json:"y"`
	IsDead bool `json:"isDead"`
}

func (w *World) Snapshot() WorldSnapshot {
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

	sumEnergy := 0.0
	for _, a := range w.Agents {
		agents = append(agents, AgentSnapshot{
			ID:     a.ID,
			X:      a.X,
			Y:      a.Y,
			IsDead: a.IsDie,
		})
		sumEnergy += float64(a.Energy)
	}

	aCopy.AvgEnergy = sumEnergy / float64(len(w.Agents))

	aCopy.Agents = agents
	if aCopy.Agents == nil {
		aCopy.Agents = []AgentSnapshot{}
	}

	return aCopy
}

func (w *World) OutOfBound(x, y int) bool {
	if x < 0 || x >= w.Width || y < 0 || y >= w.Height {
		return true
	}
	return false
}

type Chord struct {
	x, y int
}
