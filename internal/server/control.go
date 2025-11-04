package server

import (
	"net/http"
	"time"

	"github.com/ikhwanal/tinyworlds/internal/world"
)

func pauseHandler(service *world.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.TogglePause()
	}
}

func speedUpHandler(service *world.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cur := service.Interval

		if cur > 250*time.Millisecond {
			newTick := cur / 2
			service.ChangeSpeed(newTick)
		}
	}
}

func speedDownHandler(service *world.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cur := service.Interval

		if cur > 3*time.Second {
			newTick := cur * 2
			service.ChangeSpeed(newTick)
		}
	}
}

func spawnAgentHandler(service *world.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.SpawnAgent()
	}
}

func spawnFoodHandler(service *world.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.SpawnFood()
	}
}

func ControlRouter(mux *http.ServeMux, service *world.Service) {
	mux.HandleFunc("/pause", pauseHandler(service))
	mux.HandleFunc("/speed-up", speedUpHandler(service))
	mux.HandleFunc("/speed-down", speedDownHandler(service))
	mux.HandleFunc("/spawn-agent", spawnAgentHandler(service))
	mux.HandleFunc("/spawn-food", spawnFoodHandler(service))
}
