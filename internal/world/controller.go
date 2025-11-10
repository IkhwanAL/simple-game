package world

import "math/rand/v2"

type WorldController struct {
	world    *World
	CmdChan  chan Command
	stopChan chan struct{}
}

func NewWorldController(w *World) *WorldController {
	crtl := &WorldController{
		world:    w,
		CmdChan:  make(chan Command),
		stopChan: make(chan struct{}),
	}

	go crtl.loop()
	return crtl
}

func (c *WorldController) loop() {
	for {
		select {
		case cmd := <-c.CmdChan:
			switch msg := cmd.(type) {
			case CmdTick:
				c.world.Tick()
			case CmdSpawnAgent:
				agent := NewAgent(rand.IntN(c.world.Width-1), rand.IntN(c.world.Height-1), StartingEnergy)
				c.world.AddAgent(agent)
			case CmdSpawnFood:
				c.world.SpawnFood()
			case CmdSnapshot:
				msg.Reply <- c.world.Snapshot()
			case CmdStop:
				close(c.stopChan)
				return
			}
		case <-c.stopChan:
			return
		}
	}
}

func (c *WorldController) Stop() {
	c.stopChan <- struct{}{}
}
