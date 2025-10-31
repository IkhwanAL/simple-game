package world

import (
	"testing"
	"time"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

func TestAgentMovement(t *testing.T) {
	w := world.NewWorld(20, 20)
	agent := world.NewAgent(1, 5, 5, 10)

	w.AddAgent(agent)

	w.Tick()

	if agent.X == 5 && agent.Y == 5 {
		t.Errorf("Agent is Not Moving")
	}
}

func TestFoodSpawn(t *testing.T) {
	w := world.NewWorld(20, 20)

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

// How To Test The Tick?
func TestWorldService(t *testing.T) {
	service := world.NewService(world.NewWorld(20, 20))
	service.StartTick(500 * time.Millisecond)
	if len(service.Snapshot().Agents) == 0 {
		t.Fatalf("expected agents after tick")
	}
}
