package world

import (
	"testing"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

func TestAgentMovement(t *testing.T) {
	w := world.NewWorld()
	agent := world.NewAgent(5, 5)

	w.AddAgent(agent)

	w.Tick()

	if agent.X == 5 && agent.Y == 5 {
		t.Errorf("Agent is Not Moving")
	}
}

func TestFoodSpawn(t *testing.T) {
	w := world.NewWorld()

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

func TestWorldService(t *testing.T) {
	service := world.NewService(world.NewWorld())
	service.Tick()
	if len(service.Snapshot().Agents) == 0 {
		t.Fatalf("expected agents after tick")
	}
}
