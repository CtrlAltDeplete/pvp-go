package main

import (
	"database/sql"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"sort"
	"sync"
	"time"
)

type Cup struct {
	name           string
	pokemon        []PokemonDto
	moveSets       map[int64]MoveSetDto
	ids            []int64
	battleMatrix   map[int64]map[int64]float64
	pageRankMatrix *mat.Dense
	mutex          sync.Mutex
	wg             sync.WaitGroup
	current        float64
	goal           float64
	startTime      time.Time
}

func (cup *Cup) SimulateNewMoveSet(moveSetId int64) {

}

func (cup *Cup) FillBattleMatrix() {
	cup.battleMatrix = map[int64]map[int64]float64{}
	ids := make(chan int, len(cup.ids))
	for w := 0; w < 40; w++ {
		go cup.fillBattleMatrixWorker(ids)
	}
	cup.goal = float64(len(cup.ids))
	cup.current = 0.0

	for i := 0; i < len(cup.ids); i++ {
		cup.wg.Add(1)
		ids <- i
	}
	close(ids)
	cup.wg.Wait()
}

func (cup *Cup) fillBattleMatrixWorker(ids <-chan int) {
	for i := range ids {
		ally := cup.ids[i]
		battleSims := BATTLE_SIMS_DAO.FindMatchupsForAlly(ally, cup.ids)
		battleMiniMatrix := map[int64]float64{}
		for _, sim := range battleSims {
			battleMiniMatrix[sim.EnemyId()] = sim.Score()
		}
		for _, enemyId := range cup.ids {
			if _, found := battleMiniMatrix[enemyId]; found {
				continue
			} else {
				fmt.Printf("battle_simulation not found for %d vs %d\n", ally, enemyId)
				_, allyMoveSet := MOVE_SETS_DAO.FindSingleWhere("id = ?", ally)
				_, allyPokemon := POKEMON_DAO.FindById(allyMoveSet.PokemonId())
				_, allyFast := MOVES_DAO.FindById(allyMoveSet.FastMoveId())
				_, allyPrimary := MOVES_DAO.FindById(allyMoveSet.PrimaryChargeMoveId())
				allyCharges := []MoveDto{*allyPrimary}
				if allyMoveSet.SecondaryChargeMoveId() != nil {
					_, allySecondary := MOVES_DAO.FindById(*allyMoveSet.SecondaryChargeMoveId())
					allyCharges = append(allyCharges, *allySecondary)
				}
				ally := *NewPokemon(*allyPokemon, *allyFast, allyCharges)

				_, enemyMoveSet := MOVE_SETS_DAO.FindSingleWhere("id = ?", enemyId)
				_, enemyPokemon := POKEMON_DAO.FindById(enemyMoveSet.PokemonId())
				_, enemyFast := MOVES_DAO.FindById(enemyMoveSet.FastMoveId())
				_, enemyPrimary := MOVES_DAO.FindById(enemyMoveSet.PrimaryChargeMoveId())
				enemyCharges := []MoveDto{*enemyPrimary}
				if enemyMoveSet.SecondaryChargeMoveId() != nil {
					_, enemySecondary := MOVES_DAO.FindById(*enemyMoveSet.SecondaryChargeMoveId())
					enemyCharges = append(enemyCharges, *enemySecondary)
				}
				enemy := *NewPokemon(*enemyPokemon, *enemyFast, enemyCharges)

				results := DoAllBattles([]Pokemon{ally, enemy})
				allyResults := []int64{}
				enemyResults := []int64{}
				for _, result := range results {
					allyResults = append(allyResults, result[0])
					enemyResults = append(enemyResults, result[1])
				}

				params := []int64{allyMoveSet.Id(), enemyMoveSet.Id()}
				params = append(params, allyResults...)
				params = append(params, enemyMoveSet.Id(), allyMoveSet.Id())
				params = append(params, enemyResults...)

				BATTLE_SIMS_DAO.BatchCreate(params)
				battleMiniMatrix[enemyId] = float64(CalculateTotalScore([]int64{
					allyResults[0], allyResults[1], allyResults[2],
					allyResults[3], allyResults[4], allyResults[5],
					allyResults[6], allyResults[7], allyResults[8],
				}))
			}
		}
		cup.mutex.Lock()
		cup.battleMatrix[ally] = battleMiniMatrix
		cup.current += 1.0

		if i%40 == 0 || i+1 == len(cup.ids) {
			ratio := cup.current / cup.goal
			past := float64(time.Now().Sub(cup.startTime))
			totalTime := time.Duration(past / ratio)
			eta := cup.startTime.Add(totalTime)
			fmt.Printf("%f%% Finished:\tETA %s\n", cup.current*100.0/cup.goal, eta)
		}

		cup.wg.Done()
		cup.mutex.Unlock()
	}
}

func (cup *Cup) CalculateMeta() {
	cup.FillBattleMatrix()
	tmpRankings := cup.subMetaCalculation()
	total := len(tmpRankings)

	// Siphon off lower 5% and start next set
	var rankings []Ranking
	var currentMax = 0.0
	var currentMin = 0.0
	var fivePercent = total / 20
	var boost float64
	for i := 1; i < 20; i++ {
		boost = currentMax - currentMin
		for j := range tmpRankings[:fivePercent] {
			rankings = append(rankings, Ranking{tmpRankings[j].moveSet, tmpRankings[j].score + boost, nil})
		}
		if rankings == nil {
			log.Fatalf("Rankings should not be nil.")
		}
		currentMax = rankings[len(rankings)-1].score
		cup.ids = []int64{}
		for _, r := range tmpRankings[fivePercent:] {
			cup.ids = append(cup.ids, r.moveSet.Id())
		}
		tmpRankings = cup.subMetaCalculation()
		currentMin = tmpRankings[0].score
	}

	boost = currentMax - currentMin
	for i := range tmpRankings {
		rankings = append(rankings, Ranking{tmpRankings[i].moveSet, tmpRankings[i].score + boost, nil})
	}
	if rankings == nil {
		log.Fatalf("Rankings should not be nil.")
	}
	finalMin := rankings[0].score
	finalMax := rankings[len(rankings)-1].score
	sort.Slice(rankings, func(i, j int) bool {
		return rankings[i].score > rankings[j].score
	})
	pokemonRankings := map[int64][]Ranking{}
	var currentPokemonRank int64 = 1
	for i := range rankings {
		ranking := rankings[i]
		ranking.score = (ranking.score - finalMin) * 100.0 / (finalMax - finalMin)
		if pokemonRankings[ranking.moveSet.PokemonId()] == nil {
			pokemonRankings[ranking.moveSet.PokemonId()] = []Ranking{}
			ranking.pokemonRank = sql.NullInt64{currentPokemonRank, true}
			currentPokemonRank++
		} else {
			ranking.pokemonRank = sql.NullInt64{0, false}
		}
		pokemonRankings[ranking.moveSet.PokemonId()] = append(pokemonRankings[ranking.moveSet.PokemonId()], ranking)
		rankings[i] = ranking
	}
	for _, ranking := range rankings {
		err, _ := RANKINGS_DAO.Create(cup.name, ranking.moveSet.PokemonId(), ranking.moveSet.Id(), ranking.pokemonRank, ranking.score)
		CheckError(err)
	}
}

func (cup *Cup) subMetaCalculation() []Ranking {
	cup.pageRankMatrix = mat.NewDense(len(cup.ids), len(cup.ids), nil)
	ids := make(chan int, len(cup.ids))
	for w := 0; w < 10; w++ {
		go cup.calculateMetaWorker(ids)
	}
	for i := 0; i < len(cup.ids); i++ {
		cup.wg.Add(1)
		ids <- i
	}
	close(ids)
	cup.wg.Wait()
	for col := 0; col < len(cup.ids); col++ {
		colSum := 0.0
		for row := 0; row < len(cup.ids); row++ {
			colSum += cup.pageRankMatrix.At(row, col)
		}
		for row := 0; row < len(cup.ids); row++ {
			cup.pageRankMatrix.Set(row, col, cup.pageRankMatrix.At(row, col)/colSum)
		}
	}
	var data []float64
	for i := 0; i < len(cup.ids); i++ {
		data = append(data, 1.0/float64(len(cup.ids)))
	}
	var controlVector = mat.NewDense(len(cup.ids), 1, data)
	oldOrder := cup.getRankings(controlVector)
	constantRankCounter := 0
	for i := 0; i < 500; i++ {
		controlVector.Product(cup.pageRankMatrix, controlVector)
		newOrder := cup.getRankings(controlVector)
		different := false
		for j := range newOrder {
			if oldOrder[j].moveSet != newOrder[j].moveSet {
				different = true
				break
			}
		}
		oldOrder = newOrder
		if different {
			constantRankCounter++
			if constantRankCounter > 10 && i > 50 {
				break
			}
		} else {
			constantRankCounter = 0
		}
	}
	tmpRankings := cup.getRankings(controlVector)
	return tmpRankings
}

func (cup *Cup) calculateMetaWorker(rows <-chan int) {
	for row := range rows {
		var rowData []float64
		for _, enemy := range cup.ids {
			rowData = append(rowData, cup.battleMatrix[cup.ids[row]][enemy])
		}
		cup.pageRankMatrix.SetRow(row, rowData)
		cup.wg.Done()
	}
}

func (cup *Cup) getRankings(controlVector *mat.Dense) []Ranking {
	var rankings = []Ranking{}
	for i := 0; i < len(cup.ids); i++ {
		rankings = append(rankings, Ranking{cup.moveSets[cup.ids[i]], controlVector.At(i, 0), nil})
	}
	sort.Slice(rankings, func(i, j int) bool {
		if rankings[i].moveSet.PokemonId() == rankings[j].moveSet.PokemonId() {
			if 100.0*(rankings[i].score-rankings[j].score)/rankings[i].score < 1 {
				if rankings[i].moveSet.SecondaryChargeMoveId() == nil {
					return true
				} else if rankings[j].moveSet.SecondaryChargeMoveId() == nil {
					return false
				}
			}
		}
		return rankings[i].score < rankings[j].score
	})
	for i := 0; i < len(cup.ids); i++ {
		rankings[i] = Ranking{rankings[i].moveSet, float64(i) * 100.0 / float64(len(cup.ids)), nil}
	}
	return rankings
}

func (cup *Cup) CalculateOtherTables() {
	var pokemonRankings [][]int64
	var teams = make(chan int)
	var goodMatchUps = make(chan int)
	var badMatchUps = make(chan int)

	for _, ranking := range RANKINGS_DAO.FindWhere("cup = ? AND pokemon_rank IS NOT NULL", cup.name) {
		pokemonRankings = append(pokemonRankings, []int64{ranking.PokemonId(), ranking.MoveSetId(), int64(math.Round(ranking.MoveSetRank()))})
	}

	for w := 0; w < 3; w++ {
		go cup.CalculateTeams(teams, pokemonRankings)
		go cup.CalculateGoodMatchUps(goodMatchUps, pokemonRankings)
		go cup.CalculateBadMatchUps(badMatchUps, pokemonRankings)
	}
	cup.current = 0.0

	for i := range pokemonRankings {
		cup.wg.Add(3)
		teams <- i
		goodMatchUps <- i
		badMatchUps <- i
	}
	close(teams)
	close(goodMatchUps)
	close(badMatchUps)
	cup.wg.Wait()
}

func (cup *Cup) CalculateTeams(indices <-chan int, pokemonRankings [][]int64) {
	for index := range indices {
		bestScore := 0.0
		var bestTeam []int64
		pokemonId := pokemonRankings[index][0]
		moveSetId := pokemonRankings[index][1]
		for i, allyOne := range pokemonRankings {
			allyOneId := allyOne[0]
			allyOneMoveSetId := allyOne[1]
			for _, allyTwo := range pokemonRankings[i : len(pokemonRankings)-1] {
				teamScore := 0.0
				allyTwoId := allyTwo[0]
				allyTwoMoveSetId := allyTwo[1]
				for _, enemy := range pokemonRankings {
					enemyMoveSetId := enemy[1]
					enemyScore := cup.battleMatrix[moveSetId][enemyMoveSetId]
					enemyScore = math.Max(enemyScore, cup.battleMatrix[allyOneMoveSetId][enemyMoveSetId])
					enemyScore = math.Max(enemyScore, cup.battleMatrix[allyTwoMoveSetId][enemyMoveSetId])
					teamScore += enemyScore * float64(enemy[2]) / 100.0
				}
				if teamScore > bestScore {
					bestScore = teamScore
					bestTeam = []int64{pokemonId, allyOneId, allyTwoId}
				}
			}
		}
		if bestTeam == nil {
			log.Fatal("Cannot have nil bestTeam")
		}
		cup.mutex.Lock()
		err, _ := TEAM_RANKINGS_DAO.Create(cup.name, bestTeam[0], bestTeam[1], bestTeam[2], bestScore)
		CheckError(err)
		cup.wg.Done()
		cup.current++
		fmt.Printf("%f%% done\n", cup.current*100.0/float64(3*len(pokemonRankings)))
		cup.mutex.Unlock()
	}
}

func (cup *Cup) CalculateGoodMatchUps(indices <-chan int, pokemonRankings [][]int64) {
	var goodMatchUps [][]int64
	for index := range indices {
		goodMatchUps = make([][]int64, len(pokemonRankings))
		copy(goodMatchUps, pokemonRankings)
		pokemonId := pokemonRankings[index][0]
		moveSetId := pokemonRankings[index][1]
		sort.Slice(goodMatchUps, func(i, j int) bool {
			return cup.battleMatrix[moveSetId][goodMatchUps[i][1]]*float64(goodMatchUps[i][2]) > cup.battleMatrix[moveSetId][goodMatchUps[j][1]]*float64(goodMatchUps[j][2])
		})
		cup.mutex.Lock()
		err, _ := MATCH_UPS_DAO.Create(cup.name, "good", pokemonId, goodMatchUps[0][0],
			goodMatchUps[1][0], goodMatchUps[2][0])
		CheckError(err)
		cup.wg.Done()
		cup.current++
		fmt.Printf("%f%% done\n", cup.current*100.0/float64(3*len(pokemonRankings)))
		cup.mutex.Unlock()
	}
}

func (cup *Cup) CalculateBadMatchUps(indices <-chan int, pokemonRankings [][]int64) {
	var badMatchUps [][]int64
	for index := range indices {
		badMatchUps = make([][]int64, len(pokemonRankings))
		copy(badMatchUps, pokemonRankings)
		pokemonId := pokemonRankings[index][0]
		moveSetId := pokemonRankings[index][1]
		sort.Slice(badMatchUps, func(i, j int) bool {
			return cup.battleMatrix[moveSetId][badMatchUps[i][1]]/float64(badMatchUps[i][2]) < cup.battleMatrix[moveSetId][badMatchUps[j][1]]/float64(badMatchUps[j][2])
		})
		cup.mutex.Lock()
		err, _ := MATCH_UPS_DAO.Create(cup.name, "bad", pokemonId, badMatchUps[0][0],
			badMatchUps[1][0], badMatchUps[2][0])
		CheckError(err)
		cup.wg.Done()
		cup.current++
		fmt.Printf("%f%% done\n", cup.current*100.0/float64(3*len(pokemonRankings)))
		cup.mutex.Unlock()
	}
}

type Ranking struct {
	moveSet     MoveSetDto
	score       float64
	pokemonRank interface{}
}

func NewCup(name string) *Cup {
	var cup = Cup{}
	cup.name = name
	cup.pokemon = POKEMON_DAO.FindWhereInCup(name)
	cup.moveSets = map[int64]MoveSetDto{}
	cup.ids = []int64{}
	for _, pokemon := range cup.pokemon {
		moveSets := MOVE_SETS_DAO.FindWhere("pokemon_id = ? AND simulated", pokemon.Id())
		for _, moveSet := range moveSets {
			cup.moveSets[moveSet.Id()] = moveSet
			cup.ids = append(cup.ids, moveSet.Id())
		}
	}
	cup.mutex = sync.Mutex{}
	cup.wg = sync.WaitGroup{}
	cup.startTime = time.Now()
	return &cup
}
