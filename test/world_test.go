package world

import (
	"testing"
	"time"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

func TestAgentMovement(t *testing.T) {
	w := world.NewWorld(20, 20, 0)
	agent := world.NewAgent(1, 5, 5, 10)

	w.AddAgent(agent)

	w.Tick()

	if agent.X == 5 || agent.Y == 5 {
		t.Errorf("Agent is Not Moving")
	}
}

func TestFoodSpawn(t *testing.T) {
	w := world.NewWorld(20, 20, 0)

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
	w := world.NewWorld(20, 20, 0)
	a := world.NewAgent(1, 0, 0, 10)
	w.AddAgent(a)

	a.Energy = 4

	agent := a.Reproduction(2, w.Width, w.Height)

	if agent == nil {
		t.Errorf("Agent Is Not Created")
	}
}

func TestDieMechanismAsync(t *testing.T) {
	w := world.NewWorld(20, 20, 0)
	a := world.NewAgent(1, 0, 0, 10)
	w.AddAgent(a)

	a.Energy = 0
	a.Die(w, 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	if a.IsDie != true {
		t.Errorf("Agent Not Die Despite Energy Reach 0")
	}

	agents := w.Snapshot().Agents

	if len(agents) != 0 {
		t.Errorf("Agent is Not Dying")
	}
}

func TestDieMechanism(t *testing.T) {
	w := world.NewWorld(20, 20, 0)
	a := world.NewAgent(1, 0, 0, 10)
	w.AddAgent(a)

	w.RemoveAgentNow(a.ID)

	if len(w.Agents) != 0 {
		t.Errorf("Agent Is Still Alive")
	}
}
