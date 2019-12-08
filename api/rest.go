package api

import (
	"PvP-Go/db/daos"
	"encoding/json"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func getRanking(w http.ResponseWriter, r *http.Request) {
	cup := mux.Vars(r)["cup"]
	pokemonIdString := mux.Vars(r)["pokemonId"]
	pokemonId, err := strconv.Atoi(pokemonIdString)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Bad params: [pokemonId] %s", pokemonIdString)
	}
	response := daos.API_DAO.GetRanking(cup, int64(pokemonId))
	if response == nil {
		_, _ = fmt.Fprintf(w, "Invalid [cup] or [pokemonId]: %s, %s", cup, pokemonIdString)
	}
	_ = json.NewEncoder(w).Encode(response)
}

func getAllRankingsForCup(w http.ResponseWriter, r *http.Request) {
	cup := mux.Vars(r)["cup"]
	response := daos.API_DAO.GetAllRankingsForCup(cup)
	if response == nil {
		_, _ = fmt.Fprintf(w, "Invalid [cup]: %s", cup)
	}
	_ = json.NewEncoder(w).Encode(response)
}

func Rest() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/getRanking/{cup}/{pokemonId}", getRanking).Methods("GET")
	router.HandleFunc("/getCupRankings/{cup}", getAllRankingsForCup).Methods("GET")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
