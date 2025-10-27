package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
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

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt)

	go func() {
		<-sigs
		log.Print("Server Down")
		// Snapshot Occur Here
		os.Exit(1)
	}()

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

		worldSnapshot := w.Snapshot()
		ui.WorldView(&worldSnapshot).Render(r.Context(), write)
	})

	http.HandleFunc("/tick", func(write http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		worldSnapshot := w.Snapshot()

		ui.WorldView(&worldSnapshot).Render(r.Context(), write)
	})

	http.HandleFunc("/metrics", func(write http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Fprintf(write, "Alloc = %v KB\nNumGoroutine = %v\n", m.Alloc/1024, runtime.NumGoroutine())
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe("127.0.0.1:8000", nil)
}
