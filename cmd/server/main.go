package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ikhwanal/tinyworlds/internal/server"
	"github.com/ikhwanal/tinyworlds/internal/world"
)

var (
	w = world.NewWorld()
)

func main() {
	world.InitLogger()

	agent := world.NewAgent(10, 10)
	w.Agents = append(w.Agents, agent)
	w.Grid[10][10].Type = world.AgentEn

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		log.Print("Server Down")

		snapshot := w.Snapshot()

		world.StoreQuickLog("log/snapshot.log", world.ToJSONBytes(snapshot))

		os.Exit(1)
	}()

	svc := world.NewService(w)
	svc.StartTick(500 * time.Millisecond)
	server.Start(svc)
}
