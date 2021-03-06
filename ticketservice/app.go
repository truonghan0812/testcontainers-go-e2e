package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Ticket struct {
	ID    string `json:"id"`
	Movie string `json:"movie"`
	User  struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	} `json:"user"`
}

type ticketDB struct {
	UserID int    `json:"user_id"`
	Movie  string `json:"movie"`
}

var (
	inMemoryDB         = sync.Map{}
	userServiceBaseURL = os.Getenv("USER_SERVICE_URL")
	port               = os.Getenv("PORT")
)

func main() {
	r := mux.NewRouter().StrictSlash(true)
	r.NewRoute().
		Path("/health").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Write([]byte(`{"running": true}`))
			w.WriteHeader(http.StatusOK)
		})
	r.NewRoute().
		Path("/tickets").
		Methods(http.MethodPost).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body ticketDB
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			id := uuid.New().String()
			inMemoryDB.Store(id, body)
			w.WriteHeader(http.StatusCreated)

			fmt.Fprintf(w, `{"id": "%s"}`, id)

		})
	r.NewRoute().
		Path("/tickets/{id}").
		Methods(http.MethodGet).
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := mux.Vars(r)["id"]
			load, ok := inMemoryDB.Load(id)
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			ticket := load.(ticketDB)

			resp, err := http.Get(fmt.Sprintf("%s/users/%d", userServiceBaseURL, ticket.UserID))
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			var user struct {
				ID   int64  `json:"id"`
				Name string `json:"name"`
			}
			json.NewDecoder(resp.Body).Decode(&user)

			json.NewEncoder(w).Encode(Ticket{
				ID:    id,
				Movie: ticket.Movie,
				User:  user,
			})
		})
	fmt.Println("running service")

	http.ListenAndServe(":"+port, r)

}
