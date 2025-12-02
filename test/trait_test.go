package test

import (
	"testing"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

func TestTraitPushBehavior(t *testing.T) {
	a := world.NewAgent(0, 0, 100, "")

	a.Greed.Current += 40

	a.Greed.PushTowardBase(1)

	if a.Greed.Current != a.Greed.Base {
		t.Errorf("this agent greed is not calmed")
	}

	a.Greed.Current = a.Greed.Base
	a.Greed.Current += 50

	oldGreed := a.Greed.Current
	a.Greed.PushTowardBase(0.5)

	expectGreed := oldGreed / 2

	if a.Greed.Current == expectGreed {
		t.Errorf("this agent greed should be calmed in halved value")
	}
}
