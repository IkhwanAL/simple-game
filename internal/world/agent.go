package world

import (
	"container/list"
	"errors"
	"fmt"
	"math/rand/v2"
	"slices"
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

type Trait struct {
	Base    uint8
	Current uint8
}

type Agent struct {
	ID            int
	X, Y          int
	Energy        int
	IsDie         bool
	Dir           [2]int
	Path          []Chord
	Greed         Trait
	Curios        Trait
	Lazy          Trait
	FieldOfVision int
	Color         string
}

func newAgentID() int {
	return int(nextAgentId.Add(1))
}

var maxTraitValue uint8 = (1 << 8) - 1

var baseValue uint8 = 128

func NewAgent(x, y, energy int) *Agent {
	greed := uint8(rand.IntN(int(maxTraitValue)))
	curious := uint8(rand.IntN(int(maxTraitValue)))
	lazy := uint8(rand.IntN(int(maxTraitValue)))

	direction := [][2]int{UP, DOWN, LEFT, RIGHT}
	startingDirection := direction[rand.IntN(len(direction))]
	id := newAgentID()

	color := fmt.Sprintf("%x%x%x", greed, curious, lazy)

	return &Agent{
		ID:            id,
		X:             x,
		Y:             y,
		Energy:        energy,
		Dir:           startingDirection,
		Greed:         Trait{Base: max(baseValue, greed), Current: max(baseValue, greed)},
		Curios:        Trait{Base: max(baseValue, curious), Current: max(baseValue, curious)},
		Lazy:          Trait{Base: max(baseValue, lazy), Current: max(baseValue, lazy)},
		FieldOfVision: 3,
		Color:         color,
	}
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

func (a *Agent) GreedControl() {
	// Low on Resource
	if a.Energy < 15 {
		currentGreed := a.Greed.Current
		nextPossibleGreed := currentGreed + 10

		if nextPossibleGreed > maxTraitValue {
			return
		}
		a.Greed.Current += 10
	}
}

func (a *Agent) SniffForFood(w *World) bool {
	type Node struct {
		Chord Chord
		Dist  int
	}

	var visited map[Chord]bool = make(map[Chord]bool)
	var parent map[Chord]Chord = make(map[Chord]Chord)

	queue := list.New()

	start := Chord{x: a.X, y: a.Y}

	current := Node{Chord: start, Dist: 0}

	queue.PushBack(current) // Current Agent Position

	var target *Chord

	for queue.Len() > 0 {
		node := queue.Remove(queue.Front()).(Node)
		cur := node.Chord
		visited[cur] = true

		if node.Dist > 0 && (cur.x != a.X || cur.y != a.Y) && w.Grid[cur.y][cur.x].Type == Food {
			target = &cur
			break
		}

		for _, d := range DIRS {
			nx, ny := cur.x+d[0], cur.y+d[1]

			next := Chord{x: nx, y: ny}

			if w.OutOfBound(nx, ny) {
				continue
			}

			if visited[next] {
				continue
			}

			if w.Grid[ny][nx].Type == Obstacle || w.Grid[ny][nx].Type == AgentEn {
				continue
			}

			dist := node.Dist + 1
			// This Will Create Diamond Shape Field Of Vision
			if dist > a.FieldOfVision {
				continue
			}

			tempNode := Node{Chord: next, Dist: dist}
			parent[next] = cur
			queue.PushBack(tempNode)
		}
	}

	if target == nil {
		return false
	}

	path := []Chord{}
	t := *target

	for t != start {
		path = append(path, t)
		t = parent[t]
	}

	slices.Reverse(path)

	a.Path = path
	return true
}
