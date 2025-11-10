package world

import (
	"errors"
	"math/rand/v2"
	"sync/atomic"
)

var UP = [2]int{0, -1}
var DOWN = [2]int{0, 1}
var LEFT = [2]int{-1, 0}
var RIGHT = [2]int{1, 0}

var DIRS = [4][2]int{
	UP, DOWN, LEFT, RIGHT,
}

var nextAgentId atomic.Int64

type Agent struct {
	ID     int
	X, Y   int
	Energy int
	IsDie  bool
	Dir    [2]int
	Path   []Chord
}

func newAgentID() int {
	return int(nextAgentId.Add(1))
}

func NewAgent(x, y, energy int) *Agent {
	direction := [][2]int{UP, DOWN, LEFT, RIGHT}
	startingDirection := direction[rand.IntN(len(direction))]
	id := newAgentID()

	return &Agent{ID: id, X: x, Y: y, Energy: energy, Dir: startingDirection}
}

func (a *Agent) Eat(w *World) {
	if w.Grid[a.Y][a.X].Type == Food {
		w.AmountFood -= 1
		a.Energy += EnergyFoodGain
		a.Path = nil
	}
}

func (a *Agent) NextMove(w *World) ([2]int, error) {
	forward := a.Dir
	rigth := [2]int{-a.Dir[1], a.Dir[0]}
	left := [2]int{a.Dir[1], -a.Dir[0]}
	back := [2]int{-a.Dir[0], -a.Dir[1]}

	// Weight Decision
	candidates := [][2]int{
		forward, forward, forward,
		left, rigth,
		back,
	}

	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, c := range candidates {
		nx := a.X + c[0]
		ny := a.Y + c[1]

		if nx < 0 || nx >= w.Width || ny < 0 || ny >= w.Height {
			continue
		}

		location := w.Grid[ny][nx].Type
		if location == Obstacle || location == AgentEn {
			continue
		}

		a.Dir = c
		return [2]int{nx, ny}, nil
	}

	return [2]int{}, errors.New("trapped in void")
}

func (a *Agent) MoveAiminglessly(w *World) (int, int) {
	nextMove, err := a.NextMove(w)
	if err != nil {
		Logf("Error When Agent Try To Move %v", err)
		return a.X, a.Y
	}
	nx, ny := nextMove[0], nextMove[1]
	a.ReduceEnergy()
	return nx, ny
}

func (a *Agent) ReduceEnergy() {
	a.Energy -= EnergyPerTick
}

// Set Agent Position But Not Reflect it Into World Map
func (a *Agent) SetAgentPosition(px, py int) {
	a.X = px
	a.Y = py
}

func (a *Agent) Reproduction(w *World) *Agent {
	chance := rand.IntN(1000)

	success := 50

	if w.DebugMode {
		chance = 1
	}

	if chance < success && a.Energy >= EnergyReproduceThreshold {
		a.Energy -= EnergyReproduceCost
		directions := [][2]int{
			{0, -1},  // up
			{1, -1},  // up-right
			{1, 0},   // right
			{1, 1},   // down-right
			{0, 1},   // down
			{-1, 1},  // down-left
			{-1, 0},  // left
			{-1, -1}, // up-left
		}

		for _, d := range directions {
			nx := a.X + d[0]
			ny := a.Y + d[1]

			if nx < 0 || nx >= w.Width || ny < 0 || ny >= w.Height {
				continue
			}

			landmark := w.Grid[ny][nx].Type

			if landmark == Obstacle || landmark == AgentEn || landmark == Food {
				continue
			}

			return NewAgent(nx, ny, StartingEnergy)
		}

	}

	return nil
}

func (a *Agent) Die(w *World) {
	if a.Energy == 0 && !a.IsDie {
		a.IsDie = true
		w.DeathCount++
		w.PendingDead = append(w.PendingDead, a)
	}
}
