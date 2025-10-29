package world

import (
	"testing"
)

func TestAgentMovement(t *testing.T) {
	world := NewWorld()
	agent := NewAgent(5, 5)

	world.AddAgent(agent)

	world.Tick()

	if agent.X == 5 && agent.Y == 5 {
		t.Errorf("Agent is Not Moving")
	}
}

func TestFoodSpawn(t *testing.T) {
	world := NewWorld()

	for range 5 {
		world.Tick()
	}

	countFood := world.AmountFood
	if countFood == 0 {
		t.Errorf("Food Not Spawn At All")
	}

	if countFood < 0 {
		t.Errorf("Food Reach Negative Value")
	}
}
