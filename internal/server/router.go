package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
	"time"

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
		log.Printf("Incoming websocket request from %s", r.RemoteAddr)
		conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
			InsecureSkipVerify: true,
		})
		if err != nil {
			log.Println("failed to open socket")
			return
		}
		defer conn.CloseNow()
		log.Println("WebSocket accepted")

		hub.AddConn(conn)
		defer func() {
			hub.RemoveConn(conn)
			_ = conn.Close(websocket.StatusNormalClosure, "")
		}()

		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)

			_, info, err := conn.Read(ctx)
			cancel()
			if err != nil {
				code := websocket.CloseStatus(err)
				if code == websocket.StatusGoingAway || code == websocket.StatusNormalClosure {
					break
				}
				log.Printf("websocket read error: %v (code=%v)", err, code)
				return
			}

			var received struct {
				Type string
			}

			err = json.Unmarshal(info, &received)
			if err != nil {
				log.Printf("failed to receive information from client\n")
				return
			}

			if received.Type == "ping" {
				continue
			}
		}
	})

	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("assets/js"))))
	mux.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("assets/css"))))

	return mux
}
