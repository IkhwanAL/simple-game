package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ikhwanal/tinyworlds/ui"
	"github.com/ikhwanal/tinyworlds/world"
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

	http.HandleFunc("/", func(write http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		ui.WorldView(w).Render(r.Context(), write)
	})

	http.HandleFunc("/tick", func(write http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		fmt.Println("Called Tick")

		ui.WorldView(w).Render(r.Context(), write)
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe("127.0.0.1:8000", nil)
}
