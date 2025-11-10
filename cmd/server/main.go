package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ikhwanal/tinyworlds/internal/server"
	"github.com/ikhwanal/tinyworlds/internal/world"
)

func main() {

	w := world.NewWorld(20, 20, 1, false)

	world.InitLogger()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	hubSocket := server.NewWebSocketHub()

	svc := server.NewService(w)
	svc.StartTick(500*time.Millisecond, hubSocket)

	srv := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: server.Router(svc, hubSocket), // ✅ We expose a router instead of blocking
	}

	go func() {
		log.Println("Tiny World Start")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Http Error %s", err.Error())
		}
	}()

	<-ctx.Done()
	log.Println("Shutting Down...")
	svc.Stop()

	shutDownCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	srv.Shutdown(shutDownCtx)

	log.Println("Complete Shutdown")
}
