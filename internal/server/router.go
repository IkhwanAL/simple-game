package server

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"

	"github.com/coder/websocket"
	"github.com/ikhwanal/tinyworlds/internal/world"
	ui "github.com/ikhwanal/tinyworlds/templates"
)

func Router(svc *Service, hub *WebSocketHub) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(write http.ResponseWriter, r *http.Request) {
		worldSnapshot := svc.Snapshot()

		worldComp := ui.WorldView(worldSnapshot)
		err := ui.MainView(worldComp).Render(r.Context(), write)
		if err != nil {
			world.Logf("failed to return html page %v", err)
		}
	})

	mux.HandleFunc("POST /world-fragment", func(write http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(write, "failed to read body", http.StatusBadRequest)
			return
		}

		var snap world.WorldSnapshot
		err = json.Unmarshal(body, &snap)
		if err != nil {
			http.Error(write, "failed to read body", http.StatusBadRequest)
			return
		}

		err = ui.WorldBoardView(snap).Render(r.Context(), write)
		if err != nil {
			world.Logf("failed to return html page %v", err)
		}
	})

	mux.HandleFunc("GET /metrics", func(write http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Fprintf(write, "Alloc = %v KB\nNumGoroutine = %v\n", m.Alloc/1024, runtime.NumGoroutine())
	})

	ControlRouter(mux, svc)

	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Println("failed to open socket")
			return
		}

		hub.AddConn(conn)
		defer func() {
			hub.RemoveConn(conn)
			conn.Close(websocket.StatusNormalClosure, "")
		}()

		ctx := context.Background()
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				code := websocket.CloseStatus(err)
				if code == websocket.StatusGoingAway || code == websocket.StatusNormalClosure {
					// expected client close — ignore quietly
					break
				}
				log.Printf("websocket read error: %v (code=%v)", err, code)
				return
			}
		}

	})

	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("GET /js/", http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js"))))

	return mux
}
