package server

import (
	"net/http"
	"time"
)

func pauseHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.TogglePause()
	}
}

func speedUpHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cur := service.Interval

		if cur > 250*time.Millisecond {
			newTick := cur / 2
			service.ChangeSpeed(newTick)
		}
	}
}

func speedDownHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cur := service.Interval

		if cur > 3*time.Second {
			newTick := cur * 2
			service.ChangeSpeed(newTick)
		}
	}
}

func spawnAgentHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.SpawnAgent()
	}
}

func spawnFoodHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		service.SpawnFood()
	}
}

func ControlRouter(mux *http.ServeMux, service *Service) {
	mux.HandleFunc("POST /pause", pauseHandler(service))
	mux.HandleFunc("POST /speed-up", speedUpHandler(service))
	mux.HandleFunc("POST /speed-down", speedDownHandler(service))
	mux.HandleFunc("POST /spawn-agent", spawnAgentHandler(service))
	mux.HandleFunc("POST /spawn-food", spawnFoodHandler(service))
}
