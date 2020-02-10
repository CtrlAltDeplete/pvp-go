package main

import (
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
	response := API_DAO.GetRanking(cup, int64(pokemonId))
	if response.Name == "" {
		_, _ = fmt.Fprintf(w, "Invalid [cup] or [pokemonId]: %s, %s", cup, pokemonIdString)
		return
	}
	_ = json.NewEncoder(w).Encode(response)
}

func getAllRankingsForCup(w http.ResponseWriter, r *http.Request) {
	cup := mux.Vars(r)["cup"]
	response := API_DAO.GetAllRankingsForCup(cup)
	if *response == nil {
		_, _ = fmt.Fprintf(w, "Invalid [cup]: %s", cup)
	} else {
		_ = json.NewEncoder(w).Encode(response)
	}
}

func saveCard(w http.ResponseWriter, r *http.Request) {
	cup := mux.Vars(r)["cup"]
	moveSetIdString := mux.Vars(r)["moveSetId"]
	moveSetId, err := strconv.Atoi(moveSetIdString)
	if err != nil {
		_, _ = fmt.Fprintf(w, "Bad params: [moveSetId] %s", moveSetIdString)
		return
	}
	response, err := API_DAO.SaveCard(cup, int64(moveSetId))
	if err != nil {
		_, _ = fmt.Fprintf(w, "Error generating image: %s", err.Error())
	} else {
		_, _ = w.Write(response)
	}
}

func Rest() {
	router := mux.NewRouter().StrictSlash(true)

	router.Path("/getCupRankings").Queries("cup", "{cup}").HandlerFunc(getAllRankingsForCup).Name("getCupRankings")
	router.Path("/getRanking").Queries("cup", "{cup}", "pokemonId", "{pokemonId}").HandlerFunc(getRanking).Name("getRanking")
	router.Path("/saveCard").Queries("cup", "{cup}", "moveSetId", "{moveSetId}").HandlerFunc(saveCard).Name("saveCard")

	log.Fatal(http.ListenAndServeTLS(":8080", "/etc/letsencrypt/live/www.pvp-go.com/fullchain.pem", "/etc/letsencrypt/live/www.pvp-go.com/privkey.pem", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(router)))
	//log.Fatal(http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET"}), handlers.AllowedOrigins([]string{"*"}))(router)))
}
