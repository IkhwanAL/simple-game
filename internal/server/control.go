package server

import (
	"fmt"
	"net/http"
	"time"
)

func IsMethodCorrect(w http.ResponseWriter, r *http.Request, method string) error {
	if r.Method != method {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return fmt.Errorf("incorrect method, use %s", method)
	}

	return nil
}

func pauseHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := IsMethodCorrect(w, r, "POST")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		service.TogglePause()
	}
}

func speedUpHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := IsMethodCorrect(w, r, "POST")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}

		speedUp(service)
	}
}

func speedUp(service *Service) {
	cur := service.Interval
	newTick := cur / 2

	if newTick > 200*time.Millisecond {
		service.ChangeSpeed(newTick)
	}
}

func speedDownHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := IsMethodCorrect(w, r, "POST")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		speedDown(service)
	}
}

func speedDown(service *Service) {
	cur := service.Interval
	newTick := cur * 2

	if cur < 2*time.Second {
		service.ChangeSpeed(newTick)
	}
}

func spawnAgentHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := IsMethodCorrect(w, r, "POST")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		service.SpawnAgent()
	}
}

func spawnFoodHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := IsMethodCorrect(w, r, "POST")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		service.SpawnFood(1)
	}
}

func spawnMultipleFoodHandler(service *Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := IsMethodCorrect(w, r, "POST")
		if err != nil {
			fmt.Fprint(w, err)
			return
		}
		service.SpawnFood(10)
	}
}

func ControlRouter(mux *http.ServeMux, service *Service) {
	mux.HandleFunc("/pause", pauseHandler(service))
	mux.HandleFunc("/speed-up", speedUpHandler(service))
	mux.HandleFunc("/speed-down", speedDownHandler(service))
	mux.HandleFunc("/spawn-agent", spawnAgentHandler(service))
	mux.HandleFunc("/spawn-food", spawnFoodHandler(service))
	mux.HandleFunc("/spawn-multiple-food", spawnMultipleFoodHandler(service))
}
