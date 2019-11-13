package main

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"PvP-Go/models"
	"fmt"
	"log"
	"sync"
)

var (
	mutex       = sync.Mutex{}
	waitGroup   = sync.WaitGroup{}
	finished    = 0.0
	total       = 0.0
	allMoves    map[int64]dtos.MoveDto
	allPokemon  map[int64]dtos.PokemonDto
	allMovesets []dtos.MoveSetDto
)

func worker(jobs <-chan int) {
	for i := range jobs {
		allyMovesetDto := allMovesets[i]
		allyPokeDto := allPokemon[allyMovesetDto.PokemonId()]
		allyFastMoveDto := allMoves[allyMovesetDto.FastMoveId()]
		allyChargeMoveDtos := []dtos.MoveDto{allMoves[allyMovesetDto.PrimaryChargeMoveId()]}
		if allyMovesetDto.SecondaryChargeMoveId() != nil {
			allyChargeMoveDtos = append(allyChargeMoveDtos, allMoves[*allyMovesetDto.SecondaryChargeMoveId()])
		}
		ally := *models.NewPokemon(allyPokeDto, allyFastMoveDto, allyChargeMoveDtos)

		j := i
		for j < int(total) {
			enemyMovesetDto := allMovesets[j]
			enemyPokeDto := allPokemon[enemyMovesetDto.PokemonId()]
			enemyFastMoveDto := allMoves[enemyMovesetDto.FastMoveId()]
			enemyChargeMoveDtos := []dtos.MoveDto{allMoves[enemyMovesetDto.PrimaryChargeMoveId()]}
			if enemyMovesetDto.SecondaryChargeMoveId() != nil {
				enemyChargeMoveDtos = append(enemyChargeMoveDtos, allMoves[*enemyMovesetDto.SecondaryChargeMoveId()])
			}
			enemy := *models.NewPokemon(enemyPokeDto, enemyFastMoveDto, enemyChargeMoveDtos)

			results := models.DoAllBattles([]models.Pokemon{ally, enemy})
			allyResults := []int64{}
			enemyResults := []int64{}
			for _, result := range results {
				allyResults = append(allyResults, result[0])
				enemyResults = append(enemyResults, result[1])
			}

			mutex.Lock()
			err, _ := daos.BATTLE_SIMS_DAO.Create(allyMovesetDto.Id(), enemyMovesetDto.Id(), allyResults)
			if err != nil {
				log.Printf("Battle Simulation Failed on (%d, %d): %s\n", allyMovesetDto.Id(), enemyMovesetDto.Id(), err.Error())
			}
			finished++

			if allyMovesetDto.Id() != enemyMovesetDto.Id() {
				err, _ := daos.BATTLE_SIMS_DAO.Create(enemyMovesetDto.Id(), allyMovesetDto.Id(), enemyResults)
				if err != nil {
					log.Printf("Battle Simulation Failed on (%d, %d): %s\n", enemyMovesetDto.Id(), allyMovesetDto.Id(), err.Error())
				}
			}
			finished++
			if int(finished)%1000 == 0 {
				fmt.Printf("%d%% Finished\n", int(finished*100.0/(total*total)))
			}
			waitGroup.Done()
			mutex.Unlock()
			j++
		}
	}
}

func main() {
	fmt.Println("Gathering moves...")
	allMoves = map[int64]dtos.MoveDto{}
	for _, moveDto := range daos.MOVES_DAO.FindAll() {
		allMoves[moveDto.Id()] = moveDto
	}

	fmt.Println("Gathering pokemon...")
	allPokemon = map[int64]dtos.PokemonDto{}
	for _, pokemonDto := range daos.POKEMON_DAO.FindAll() {
		allPokemon[pokemonDto.Id()] = pokemonDto
	}

	fmt.Println("Gathering move sets...")
	allMovesets = daos.MOVE_SETS_DAO.FindAll()

	fmt.Println("Preparing workers...")
	total = float64(len(allMovesets))
	jobs := make(chan int, int(total))

	for w := 0; w < 4; w++ {
		go worker(jobs)
	}

	fmt.Println("Starting work...")
	for i := range allMovesets {
		waitGroup.Add(1)
		jobs <- i
	}
	close(jobs)
	waitGroup.Wait()
}
