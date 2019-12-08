package processes

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"PvP-Go/models"
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
	allMoves      map[int64]dtos.MoveDto
	allPokemon    map[int64]dtos.PokemonDto
	allMoveSets   []dtos.MoveSetDto
	startTime     time.Time
	batchParams   []int64
)

func addToBatch(allyId, enemyId int64, allyResults []int64) {
	sqlQueueGroup.Wait()
	batchParams = append(batchParams, allyId, enemyId)
	batchParams = append(batchParams, allyResults...)
	if len(batchParams) >= 65000 {
		sqlQueueGroup.Add(1)
		var oldParams = make([]int64, len(batchParams))
		copy(oldParams, batchParams)
		sqlDoneGroup.Wait()
		go func() {
			sqlDoneGroup.Add(1)
			daos.BATTLE_SIMS_DAO.BatchCreate(oldParams)
			sqlDoneGroup.Done()
		}()
		sqlQueueGroup.Done()
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
		allyChargeMoveDtos := []dtos.MoveDto{allMoves[allyMovesetDto.PrimaryChargeMoveId()]}
		if allyMovesetDto.SecondaryChargeMoveId() != nil {
			allyChargeMoveDtos = append(allyChargeMoveDtos, allMoves[*allyMovesetDto.SecondaryChargeMoveId()])
		}
		ally := *models.NewPokemon(allyPokeDto, allyFastMoveDto, allyChargeMoveDtos)

		j := 0
		for j <= i {
			enemyMovesetDto := allMoveSets[j]
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

func Simulate() {
	finished = 0.0
	total = 0.0

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
	allMoveSets = daos.MOVE_SETS_DAO.FindWhere("1 = 1 ORDER BY id ASC")

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

	for w := 0; w < 12; w++ {
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
		daos.BATTLE_SIMS_DAO.BatchCreate(batchParams)
	}
}
