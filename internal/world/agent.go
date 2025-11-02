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
		a.Energy += 10
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

func (a *Agent) Move(w *World) {
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

func (a *Agent) Reproduction(ID int, w *World) *Agent {
	thresholdEnergy := 8

	if a.Energy > thresholdEnergy {
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

			if landmark == Obstacle || landmark == AgentEn {
				continue
			}

			return NewAgent(ID+1, nx, ny, StartingEnergy)
		}

	}

	return nil
}

func (a *Agent) Die(w *World, dieDuration time.Duration) {
	if a.Energy == 0 {
		w.RemoveAgent(a, dieDuration)
	}
}
