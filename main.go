package main

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"PvP-Go/models"
	"fmt"
	"sync"
)

var (
	mutex     = sync.Mutex{}
	waitGroup = sync.WaitGroup{}
	finished  = 0.0
	total     = 0.0
	fastMove  dtos.MoveDto
)

func worker(jobs <-chan dtos.PokemonDto) {
	for pokemonDto := range jobs {
		models.NewPokemon(pokemonDto, fastMove, []dtos.MoveDto{})
		mutex.Lock()
		finished++
		fmt.Printf("%d%% Finished\n", int(finished*100.0/total))
		waitGroup.Done()
		mutex.Unlock()
	}
}

func main() {
	fastMove = daos.MOVES_DAO.FindAll()[0]
	allPokemon := daos.POKEMON_DAO.FindAll()
	total = float64(len(allPokemon))
	jobs := make(chan dtos.PokemonDto, len(allPokemon))

	for w := 0; w < 4; w++ {
		go worker(jobs)
	}

	for _, pokeDto := range allPokemon {
		waitGroup.Add(1)
		jobs <- pokeDto
	}
	close(jobs)
	waitGroup.Wait()
}
