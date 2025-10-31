package world

import (
	"errors"
	"math/rand/v2"
	"time"
)

var UP = [2]int{0, -1}
var DOWN = [2]int{0, 1}
var LEFT = [2]int{-1, 0}
var RIGHT = [2]int{1, 0}

type Agent struct {
	ID     int
	X, Y   int
	Energy int
	IsDie  bool
	Dir    [2]int
}

func NewAgent(id, x, y, energy int) *Agent {
	direction := [][2]int{UP, DOWN, LEFT, RIGHT}
	startingDirection := direction[rand.IntN(len(direction))]

	return &Agent{ID: id, X: x, Y: y, Energy: energy, Dir: startingDirection}
}

func (a *Agent) Eat(w *World) {
	if w.Grid[a.Y][a.X].Type == Food {
		w.Grid[a.Y][a.X].Type = Empty
		w.AmountFood -= 1
		a.Energy += 5
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

func (a *Agent) CheckSurround(w *World) ([2]int, error) {
	// x, y
	leftTop := [2]int{-1, -1}
	middleTop := [2]int{0, -1}
	rightTop := [2]int{1, -1}

	leftMiddle := [2]int{-1, 0}
	rightMiddle := [2]int{1, 0}

	leftBottom := [2]int{-1, 1}
	middleBottom := [2]int{0, 1}
	rightBottom := [2]int{1, 1}

	surround := [][2]int{
		leftTop,
		middleTop,
		rightTop,
		leftMiddle,
		rightMiddle,
		leftBottom,
		middleBottom,
		rightBottom,
	}

	freeLocation := make([][2]int, 0)

	for _, mark := range surround {
		x, y := mark[0], mark[1]

		nx := a.X + x
		ny := a.Y + y

		if (nx >= 0 && nx < w.Width) && (ny >= 0 && ny < w.Height) {
			location := w.Grid[ny][nx].Type

			if location == Obstacle || location == AgentEn {
				continue
			}

			freeLocation = append(freeLocation, mark)
		}

	}

	if len(freeLocation) == 0 {
		return [2]int{}, errors.New("they trapped")
	}

	rand.Shuffle(len(freeLocation), func(i, j int) {
		freeLocation[i], freeLocation[j] = freeLocation[j], freeLocation[i]
	})

	nextLocation := [2]int{a.X + freeLocation[0][0], a.Y + freeLocation[0][1]}
	return nextLocation, nil
}

func (a *Agent) Move2(w *World) {
	nextMove, err := a.NextMove(w)
	if err != nil {
		Logf("Error When Agent Try To Move %v", err)
		return
	}

	w.Grid[a.Y][a.X].Type = Empty

	nx, ny := nextMove[0], nextMove[1]

	a.X = nx
	a.Y = ny

	a.Energy--

	w.Grid[ny][nx].Type = AgentEn

}

func (a *Agent) Move(w *World) {
	dx := rand.IntN(3) - 1
	dy := rand.IntN(3) - 1

	nx := a.X + dx
	ny := a.Y + dy

	if (nx > 0 && nx < w.Width) && (ny > 0 && ny < w.Height) {

		w.Grid[a.Y][a.X].Type = Empty

		a.X = nx
		a.Y = ny

		a.Energy--

		w.Grid[ny][nx].Type = AgentEn
	}
}

func (a *Agent) Reproduction(ID, worldWidth, worldHeight int) *Agent {
	thresholdEnergy := 4

	if a.Energy < thresholdEnergy {
		nx := min(a.X+1, worldWidth)
		ny := min(a.Y+1, worldHeight)
		return NewAgent(ID+1, nx, ny, StartingEnergy)
	}

	return nil
}

func (a *Agent) Die(w *World, dieDuration time.Duration) {
	if a.Energy == 0 {
		w.RemoveAgent(a, dieDuration)
	}
}
