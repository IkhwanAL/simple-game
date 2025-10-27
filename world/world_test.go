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
