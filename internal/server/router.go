package server

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/ikhwanal/tinyworlds/internal/world"
	ui "github.com/ikhwanal/tinyworlds/templates"
)

func Router(svc *world.Service) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(write http.ResponseWriter, r *http.Request) {
		worldSnapshot := svc.Snapshot()
		err := ui.WorldView(&worldSnapshot).Render(r.Context(), write)
		if err != nil {
			world.Logf("failed to return html page %v", err)
		}
	})

	mux.HandleFunc("/tick", func(write http.ResponseWriter, r *http.Request) {
		worldSnapshot := svc.Snapshot()
		err := ui.WorldView(&worldSnapshot).Render(r.Context(), write)
		if err != nil {
			world.Logf("failed to return html page %v", err)
		}
	})

	mux.HandleFunc("/metrics", func(write http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Fprintf(write, "Alloc = %v KB\nNumGoroutine = %v\n", m.Alloc/1024, runtime.NumGoroutine())
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}
