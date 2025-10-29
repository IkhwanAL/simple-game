package main

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ikhwanal/tinyworlds/internal/server"
	"github.com/ikhwanal/tinyworlds/internal/world"
)

var (
	w  = world.NewWorld()
	mu sync.Mutex
)

func main() {
	world.InitLogger()

	agent := world.NewAgent(10, 10)
	w.Agents = append(w.Agents, agent)
	w.Grid[10][10].Type = world.AgentEn

	go func() {
		for {
			time.Sleep(1 * time.Second)
			w.Tick()
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		log.Print("Server Down")

		snapshot := w.Snapshot()

		world.StoreQuickLog("snapshot.log", world.ToJSONBytes(snapshot))

		os.Exit(1)
	}()

	svc := world.NewService(w)
	server.Start(svc)
}
