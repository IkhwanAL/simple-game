package world

import "math/rand/v2"

type Agent struct {
	X, Y int
}

func NewAgent(x, y int) *Agent {
	return &Agent{X: x, Y: y}
}

func (a *Agent) Act(w *World) {
	dx := []int{0, 1, -1, 0}
	dy := []int{1, 0, 0, -1}

	i := rand.IntN(4)
	nx := a.X + dx[i]
	ny := a.Y + dy[i]

	Logf("Nx: %d, Ny: %d", nx, ny)
	if nx >= 0 && nx < Width && ny >= 0 && ny < Height {
		if w.Grid[ny][nx].Type == Food {
			w.Grid[ny][nx].Type = Empty
			w.AmountFood -= 1
		}

		w.Grid[a.Y][a.X].Type = Empty
		a.X, a.Y = nx, ny
		w.Grid[a.Y][a.X].Type = AgentEn
	}
}
