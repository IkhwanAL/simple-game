package server

import (
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/ikhwanal/tinyworlds/internal/world"
	ui "github.com/ikhwanal/tinyworlds/templates"
)

func Start(svc *world.Service) {
	http.HandleFunc("/", func(write http.ResponseWriter, r *http.Request) {
		worldSnapshot := svc.Snapshot()
		err := ui.WorldView(&worldSnapshot).Render(r.Context(), write)
		if err != nil {
			world.Logf("failed to return html page %v", err)
		}
	})

	http.HandleFunc("/tick", func(write http.ResponseWriter, r *http.Request) {
		worldSnapshot := svc.Snapshot()
		err := ui.WorldView(&worldSnapshot).Render(r.Context(), write)
		if err != nil {
			world.Logf("failed to return html page %v", err)
		}
	})

	http.HandleFunc("/metrics", func(write http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Fprintf(write, "Alloc = %v KB\nNumGoroutine = %v\n", m.Alloc/1024, runtime.NumGoroutine())
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	err := http.ListenAndServe("127.0.0.1:8000", nil)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("TinyWorlds server running at :8000")
}
