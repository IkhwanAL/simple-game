package world

import (
	"testing"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

func TestAgentMovement(t *testing.T) {
	w := world.NewWorld(20, 20, 0, false)
	agent := world.NewAgent(5, 5, 10)

	w.AddAgent(agent)

	w.Tick()

	if agent.X == 5 && agent.Y == 5 {
		t.Errorf("Agent is Not Moving")
	}
}

func TestFoodSpawn(t *testing.T) {
	w := world.NewWorld(20, 20, 0, false)

	for range 5 {
		w.Tick()
	}

	countFood := w.AmountFood
	if countFood == 0 {
		t.Errorf("Food Not Spawn At All")
	}

	if countFood < 0 {
		t.Errorf("Food Reach Negative Value")
	}
}

func TestReproductionMechanism(t *testing.T) {
	w := world.NewWorld(20, 20, 0, true)
	a := world.NewAgent(0, 0, 10)
	w.AddAgent(a)

	a.Energy = 20

	agent := a.Reproduction(w)

	if agent == nil {
		t.Errorf("Agent Is Not Created")
	}
}

func TestDieMechanismAsync(t *testing.T) {
	w := world.NewWorld(20, 20, 0, false)
	a := world.NewAgent(0, 0, 10)
	w.AddAgent(a)

	a.Energy = 0
	a.Die(w)

	if a.IsDie != true {
		t.Errorf("Agent Not Die Despite Energy Reach 0")
	}

	w.Tick()

	agents := w.Snapshot().Agents

	if len(agents) != 0 {
		t.Errorf("Agent is Not Dying")
	}
}

func TestDieMechanism(t *testing.T) {
	w := world.NewWorld(20, 20, 0, false)
	a := world.NewAgent(0, 0, 10)
	w.AddAgent(a)

	w.RemoveAgent(a)

	if len(w.Agents) != 0 {
		t.Errorf("Agent Is Still Alive")
	}
}
