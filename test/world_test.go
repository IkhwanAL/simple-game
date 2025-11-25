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

func removeAllFood(w *world.World) {
	for y := range w.Height {
		for x := range w.Width {
			if w.Grid[y][x].Type == world.Food {
				w.Grid[y][x].Type = world.Empty
			}
		}
	}
}

func spawnOneFood(w *world.World, a *world.Agent, distanceFood int) {
	agentX, agentY := a.X, a.Y

	w.Grid[agentY+distanceFood][agentX+distanceFood].Type = world.Food
}

func TestAgentToDetectFood(t *testing.T) {
	w := world.NewWorld(20, 20, 1, true)
	removeAllFood(w)

	a := world.NewAgent(0, 0, 100)
	spawnOneFood(w, a, 1)

	w.AddAgent(a)

	// Test Detect Closest Food And Able to Found Food
	foundFood := a.SniffForFood(w)

	if !foundFood {
		t.Errorf("Agent is blind can't find a food even i already put a close food near agent")
	}

	// Test Detect Closest Food But No Food Found
	removeAllFood(w)

	foundFood = a.SniffForFood(w)

	if foundFood {
		t.Errorf("Agent Find an invisible food which is impossible")
	}
}

func TestAgentVisionFieldOfView(t *testing.T) {
	w := world.NewWorld(20, 20, 1, true)
	removeAllFood(w)

	a := world.NewAgent(0, 0, 100)

	a.FieldOfVision += 4

	spawnOneFood(w, a, 4)

	w.AddAgent(a)

	foundFood := a.SniffForFood(w)

	if foundFood {
		t.Errorf("Agent is blind need a glasess here, the vision is not working")
	}

	a.FieldOfVision -= 4

	foundFood = a.SniffForFood(w)

	if foundFood {
		t.Errorf("Agent is suppose to be blind here, the vision is not working")
	}
}
