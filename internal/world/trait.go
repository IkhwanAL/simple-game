package world

type Trait struct {
	Base    uint8
	Current uint8
}

func (t *Trait) PushTowardBase(rate float64) {
	delta := float64(t.Current-t.Base) * rate
	t.Current -= uint8(delta)
}
