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
		return
	}
	response := daos.API_DAO.GetRanking(cup, int64(pokemonId))
	if response.Name == "" {
		_, _ = fmt.Fprintf(w, "Invalid [cup] or [pokemonId]: %s, %s", cup, pokemonIdString)
		return
	}
	_ = json.NewEncoder(w).Encode(response)
}

func getAllRankingsForCup(w http.ResponseWriter, r *http.Request) {
	cup := mux.Vars(r)["cup"]
	response := daos.API_DAO.GetAllRankingsForCup(cup)
	if *response == nil {
		_, _ = fmt.Fprintf(w, "Invalid [cup]: %s", cup)
	} else {
		_ = json.NewEncoder(w).Encode(response)
	}
}

func Rest() {
	router := mux.NewRouter().StrictSlash(true)
	router.Path("/getCupRankings").Queries("cup", "{cup}").HandlerFunc(getAllRankingsForCup).Name("getCupRankings")
	router.Path("/getRanking").Queries("cup", "{cup}", "pokemonId", "{pokemonId}").HandlerFunc(getRanking).Name("getRanking")
	log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
