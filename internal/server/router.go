package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/coder/websocket"
	"github.com/ikhwanal/tinyworlds/internal/world"
	ui "github.com/ikhwanal/tinyworlds/templates"
)

func Router(svc *Service, hub *WebSocketHub) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(write http.ResponseWriter, r *http.Request) {
		worldSnapshot := svc.Snapshot()

		worldComp := ui.WorldView(worldSnapshot)
		err := ui.MainView(worldComp).Render(r.Context(), write)
		if err != nil {
			world.Logf("failed to return html page %v", err)
		}
	})

	mux.HandleFunc("/metrics", func(write http.ResponseWriter, r *http.Request) {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Fprintf(write, "Alloc = %v KB\nNumGoroutine = %v\n", m.Alloc/1024, runtime.NumGoroutine())
	})

	ControlRouter(mux, svc)

	mux.HandleFunc("/listen", func(w http.ResponseWriter, r *http.Request) {
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Println("failed to open socket")
			return
		}

		fmt.Println(conn.Subprotocol())

		hub.AddConn(conn)
		defer func() {
			hub.RemoveConn(conn)
			_ = conn.Close(websocket.StatusNormalClosure, "")
		}()

		ctx := context.Background()
		for {
			_, _, err := conn.Read(ctx)
			if err != nil {
				code := websocket.CloseStatus(err)
				if code == websocket.StatusGoingAway || code == websocket.StatusNormalClosure {
					break
				}
				log.Printf("websocket read error: %v (code=%v)", err, code)
				return
			}
		}
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js"))))

	return mux
}
