package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Item struct {
	ID   uuid.UUID `json:"_id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router
	shoppingItems []Item
}

func initServer() *Server {
	s := &Server{
		Router:        mux.NewRouter(),
		shoppingItems: []Item{},
	}
	return s
}

func addShoppingItem(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var i Item
		if err := json.NewDecoder(r.Body).Decode(&i); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		i.ID = uuid.New()
		s.shoppingItems = append(s.shoppingItems, i)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(i); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
