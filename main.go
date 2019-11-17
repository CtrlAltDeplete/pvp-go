package main

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"PvP-Go/models"
	"fmt"
	"sync"
	"time"
)

var (
	mutex        = sync.Mutex{}
	simWaitGroup = sync.WaitGroup{}
	finished     = 0.0
	total        = 0.0
	allMoves     map[int64]dtos.MoveDto
	allPokemon   map[int64]dtos.PokemonDto
	allMovesets  []dtos.MoveSetDto
	startTime    time.Time
	batchParams  []int64
)

func addToBatch(allyId, enemyId int64, allyResults []int64) {
	mutex.Lock()
	batchParams = append(batchParams, allyId, enemyId)
	batchParams = append(batchParams, allyResults...)
	if len(batchParams) >= 5000*11 {
		var oldParams = make([]int64, len(batchParams))
		copy(oldParams, batchParams)
		batchParams = []int64{}
		daos.BATTLE_SIMS_DAO.BatchCreate(oldParams)

		ratio := finished / (total * total)
		past := float64(time.Now().Sub(startTime))
		totalTime := time.Duration(past / ratio)
		eta := startTime.Add(totalTime)
		fmt.Printf("%f%% Finished:\tETA %s\n", finished*100.0/(total*total), eta)
	}
	mutex.Unlock()
}

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
			finished++
			addToBatch(allyMovesetDto.Id(), enemyMovesetDto.Id(), allyResults)
			if allyMovesetDto.Id() != enemyMovesetDto.Id() {
				finished++
				addToBatch(enemyMovesetDto.Id(), allyMovesetDto.Id(), enemyResults)
			}
			mutex.Unlock()
			j++
		}
		allyMovesetDto.SetSimulated(true)
		mutex.Lock()
		daos.MOVE_SETS_DAO.Update(allyMovesetDto)
		mutex.Unlock()
		simWaitGroup.Done()
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

	for w := 0; w < 40; w++ {
		go worker(jobs)
	}

	fmt.Println("Starting work...")
	startTime = time.Now()
	for i := range allMovesets {
		simWaitGroup.Add(1)
		jobs <- i
	}
	close(jobs)
	simWaitGroup.Wait()
	if len(batchParams) > 0 {
		daos.BATTLE_SIMS_DAO.BatchCreate(batchParams)
	}
}
