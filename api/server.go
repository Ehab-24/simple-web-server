package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

/*******************************************************
 * * * * * * * * * * Types * * * * * * * * * * * * *
 ********************************************************/
type Item struct {
	ID   uuid.UUID `json:"_id"`
	Name string    `json:"name"`
}

type Server struct {
	*mux.Router
	shoppingItems []Item
}

/*******************************************************
 * * * * * * * * * * Main * * * * * * * * * * * * *
 ********************************************************/
func NewServer() *Server {
	s := &Server{
		Router:        mux.NewRouter(),
		shoppingItems: []Item{},
	}
	s.routes()
	return s
}

/*******************************************************
 * * * * * * * * * * Routes * * * * * * * * * * * * *
 ********************************************************/
func (s *Server) routes() {
	s.HandleFunc("/", defaultHandler).Methods("GET")
	s.HandleFunc("/items", addShoppingItem(s)).Methods("POST")
	s.HandleFunc("/items", getShoppingItems(s)).Methods("GET")
	s.HandleFunc("/items/{id}", removeShoppingItem(s)).Methods("DELETE")
}

/*******************************************************
 * * * * * * * * * * Controllers * * * * * * * * * * * * *
 ********************************************************/
func defaultHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("API is up"))
}

// Create
func (i *Item) validate() bool {
	return i.Name != ""
}

func addShoppingItem(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var i Item
		err := json.NewDecoder(r.Body).Decode(&i)
		isValid := i.validate()
		if err != nil || !isValid {
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

// Read
func getShoppingItems(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(s.shoppingItems); err != nil {
			http.Error(w, "An error occured", http.StatusInternalServerError)
		}
	}
}

// Delete
func remove(slice []Item, index int) []Item {
	return append(slice[:index], slice[index+1:]...)
}

func removeShoppingItem(s *Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		idStr := mux.Vars(r)["id"]
		id, err := uuid.Parse(idStr)
		if err != nil {
			http.Error(w, "Invalid item id", http.StatusBadRequest)
			return
		}

		for i := 0; i < len(s.shoppingItems); i++ {
			if s.shoppingItems[i].ID == id {
				remove(s.shoppingItems, i)
				break
			}
		}
	}
}
