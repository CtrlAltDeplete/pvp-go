package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex         = sync.Mutex{}
	simWaitGroup  = sync.WaitGroup{}
	sqlQueueGroup = sync.WaitGroup{}
	sqlDoneGroup  = sync.WaitGroup{}
	finished      float64
	total         float64
	allMoves      map[int64]MoveDto
	allPokemon    map[int64]PokemonDto
	allMoveSets   []MoveSetDto
	startTime     time.Time
	batchParams   []int64
)

func addToBatch(allyId, enemyId int64, allyResults []int64) {
	sqlQueueGroup.Wait()
	batchParams = append(batchParams, allyId, enemyId)
	batchParams = append(batchParams, allyResults...)
	if len(batchParams) >= 65000 {
		fmt.Println("sqlQueueGroup.Add(1)")
		sqlQueueGroup.Add(1)
		defer sqlQueueGroup.Done()
		var oldParams = make([]int64, len(batchParams))
		copy(oldParams, batchParams)
		fmt.Println("sqlDoneGroup.Wait()")
		sqlDoneGroup.Wait()
		go func() {
			fmt.Println("sqlDoneGroup.Add(1)")
			sqlDoneGroup.Add(1)
			defer sqlDoneGroup.Done()
			BATTLE_SIMS_DAO.BatchCreate(oldParams)
			fmt.Println("sqlDoneGroup.Done()")
		}()
		fmt.Println("sqlQueueGroup.Done()")
		batchParams = []int64{}

		ratio := finished / total
		past := float64(time.Now().Sub(startTime))
		totalTime := time.Duration(past / ratio)
		eta := startTime.Add(totalTime)
		fmt.Printf("%f%% Finished:\tETA %s\n", finished*100.0/total, eta)
	}
}

func worker(jobs <-chan int) {
	for i := range jobs {
		allyMovesetDto := allMoveSets[i]
		allyPokeDto := allPokemon[allyMovesetDto.PokemonId()]
		allyFastMoveDto := allMoves[allyMovesetDto.FastMoveId()]
		allyChargeMoveDtos := []MoveDto{allMoves[allyMovesetDto.PrimaryChargeMoveId()]}
		if allyMovesetDto.SecondaryChargeMoveId() != nil {
			allyChargeMoveDtos = append(allyChargeMoveDtos, allMoves[*allyMovesetDto.SecondaryChargeMoveId()])
		}
		ally := *NewPokemon(allyPokeDto, allyFastMoveDto, allyChargeMoveDtos)

		j := 0
		for j <= i {
			enemyMovesetDto := allMoveSets[j]
			enemyPokeDto := allPokemon[enemyMovesetDto.PokemonId()]
			enemyFastMoveDto := allMoves[enemyMovesetDto.FastMoveId()]
			enemyChargeMoveDtos := []MoveDto{allMoves[enemyMovesetDto.PrimaryChargeMoveId()]}
			if enemyMovesetDto.SecondaryChargeMoveId() != nil {
				enemyChargeMoveDtos = append(enemyChargeMoveDtos, allMoves[*enemyMovesetDto.SecondaryChargeMoveId()])
			}
			enemy := *NewPokemon(enemyPokeDto, enemyFastMoveDto, enemyChargeMoveDtos)

			results := DoAllBattles([]Pokemon{ally, enemy})
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
				addToBatch(enemyMovesetDto.Id(), allyMovesetDto.Id(), enemyResults)
			}
			mutex.Unlock()
			j++
		}
		allyMovesetDto.SetSimulated(true)
		mutex.Lock()
		MOVE_SETS_DAO.Update(allyMovesetDto)
		mutex.Unlock()
		simWaitGroup.Done()
	}
}

func Simulate() {
	finished = 0.0
	total = 0.0

	fmt.Println("Gathering moves...")
	allMoves = map[int64]MoveDto{}
	for _, moveDto := range MOVES_DAO.FindAll() {
		allMoves[moveDto.Id()] = moveDto
	}

	fmt.Println("Gathering pokemon...")
	allPokemon = map[int64]PokemonDto{}
	for _, pokemonDto := range POKEMON_DAO.FindAll() {
		allPokemon[pokemonDto.Id()] = pokemonDto
	}

	fmt.Println("Gathering move sets...")
	allMoveSets = MOVE_SETS_DAO.FindWhere("1 = 1 ORDER BY NOT simulated")

	fmt.Println("Preparing workers...")
	jobCount := 0
	for i := range allMoveSets {
		if !allMoveSets[i].Simulated() {
			jobCount++
			total += float64(i)
		}
	}
	fmt.Printf("Found %d new move sets to simulate.\n", jobCount)
	jobs := make(chan int, jobCount)

	for w := 0; w < 24; w++ {
		go worker(jobs)
	}

	fmt.Println("Starting work...")
	startTime = time.Now()
	for i := range allMoveSets {
		if !allMoveSets[i].Simulated() {
			simWaitGroup.Add(1)
			jobs <- i
		}
	}
	close(jobs)
	sqlDoneGroup.Wait()
	simWaitGroup.Wait()
	if len(batchParams) > 0 {
		BATTLE_SIMS_DAO.BatchCreate(batchParams)
	}
}
