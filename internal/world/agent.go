package world

import (
	"math/rand/v2"
)

type Agent struct {
	ID     int
	X, Y   int
	Energy int
	IsDie  bool
}

func NewAgent(id, x, y, energy int) *Agent {
	return &Agent{ID: id, X: x, Y: y, Energy: energy}
}

func (a *Agent) Eat(w *World) {

	if w.Grid[a.Y][a.X].Type == Food {
		w.Grid[a.Y][a.X].Type = Empty
		w.AmountFood -= 1
		a.Energy += 15
	}
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
	thresholdEnergy := 10

	if a.Energy < thresholdEnergy {
		nx := min(a.X+1, worldWidth)
		ny := min(a.Y+1, worldHeight)
		return NewAgent(ID+1, nx, ny, StartingEnergy)
	}

	return nil
}

func (a *Agent) Die(w *World) {
	if a.Energy == 0 {
		w.RemoveAgent(a)
	}
}
