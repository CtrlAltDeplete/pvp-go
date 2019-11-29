package models

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"database/sql"
	"fmt"
	"gonum.org/v1/gonum/mat"
	"log"
	"math"
	"sort"
	"sync"
)

var (
	TEST     = newCup("test", "id IN (2, 3, 5, 6, 8, 9)")
	SINISTER = newCup("sinister", `id IN (
SELECT p.id
FROM pokemon p
LEFT JOIN types t ON p.type_id = t.id
WHERE (
    t.first_type IN ('ghost', 'psychic', 'steel', 'fighting')
    OR t.second_type IN ('ghost', 'psychic', 'steel', 'fighting')
) AND t.first_type != 'dark'
AND NOT t.second_type <=> 'dark'
AND p.name NOT IN ('Skarmory', 'Hypno')
AND p.is_pvp_eligible
AND NOT p.is_legendary)`)
	FEROCIOUS = newCup("ferocious", `id IN (
SELECT p.id
FROM pokemon p
WHERE p.name IN ('Absol', 'Aggron', 'Ampharos', 'Arcanine', 'Aron', 'Bagon', 'Bibarel', 'Bidoof', 'Blitzle', 'Buizel', 
'Buneary', 'Camerupt', 'Cranidos', 'Cubone', 'Delcatty', 'Donphan', 'Drillbur', 'Eevee', 'Electrike', 'Entei', 'Espeon',
'Excadrill', 'Exploud', 'Flaaffy', 'Flareon', 'Floatzel', 'Furret', 'Gabite', 'Garchomp', 'Gible', 'Girafarig', 
'Glaceon', 'Glameow', 'Granbull', 'Growlithe', 'Grumpig', 'Heatmor', 'Herdier', 'Hippopotas', 'Hippowdon', 'Houndoom', 
'Houndour', 'Jolteon', 'Kangaskhan', 'Lairon', 'Larvitar', 'Lickilicky', 'Lickitung', 'Liepard', 'Lillipup', 'Linoone', 
'Lopunny', 'Luxio', 'Luxray', 'Mamoswine', 'Manectric', 'Mareep', 'Marowak', 'Alolan Marowak', 'Meowth', 'Alolan Meowth',
'Mightyena', 'Miltank', 'Minun', 'Nidoking', 'Nidoqueen', 'Nidoran♀', 'Nidoran♂', 'Nidorina', 'Nidorino', 'Ninetales', 
'Alolan Ninetales', 'Numel', 'Pachirisu', 'Patrat', 'Persian', 'Alolan Persian', 'Phanpy', 'Pichu', 'Pikachu', 
'Piloswine', 'Plusle', 'Ponyta', 'Poochyena', 'Pupitar', 'Purrloin', 'Purugly', 'Raichu', 'Alolan Raichu', 'Raikou', 
'Rampardos', 'Rapidash', 'Raticate', 'Alolan Raticate', 'Rattata', 'Alolan Rattata', 'Rhydon', 'Rhyhorn', 'Rhyperior', 
'Sandslash', 'Alolan Sandslash', 'Sandshrew', 'Alolan Sandshrew', 'Sentret', 'Shelgon', 'Shinx', 'Skitty', 'Skuntank', 
'Smeargle', 'Sneasel', 'Spoink', 'Stantler', 'Stoutland', 'Stunky', 'Suicune', 'Swinub', 'Tauros', 'Teddiursa', 
'Torkoal', 'Tyranitar', 'Umbreon', 'Ursaring', 'Vaporeon', 'Vulpix', 'Alolan Vulpix', 'Watchog', 'Weavile', 'Zangoose', 
'Zebstrika', 'Zigzagoon'))`)
	TIMELESS = newCup("timeless", `id IN (SELECT p.id
FROM pokemon p
WHERE p.name IN ('Abomasnow', 'Absol', 'Ampharos', 'Anorith', 'Arbok', 'Arcanine', 'Ariados', 'Armaldo', 'Bagon', 
'Banette', 'Barboach', 'Bayleef', 'Beedrill', 'Bellossom', 'Bellsprout', 'Blastoise', 'Blaziken', 'Bonsly', 'Budew', 
'Buizel', 'Bulbasaur', 'Burmy', 'Cacnea', 'Cacturne', 'Camerupt', 'Carnivine', 'Carvanha', 'Cascoon', 
'Castform (Rainy)', 'Castform (Snowy)', 'Castform (Sunny)', 'Caterpie', 'Charizard', 'Charmander', 'Charmeleon', 
'Cherrim', 'Cherubi', 'Chikorita', 'Chimchar', 'Chinchou', 'Clamperl', 'Cloyster', 'Combusken', 'Corphish', 'Corsola', 
'Cradily', 'Cranidos', 'Crawdaunt', 'Croconaw', 'Cubone', 'Cyndaquil', 'Dewgong', 'Diglett', 'Donphan', 'Dragonair', 
'Drapion', 'Dratini', 'Dugtrio', 'Dusclops', 'Dusknoir', 'Duskull', 'Dustox', 'Ekans', 'Electabuzz', 'Electivire', 
'Electrike', 'Electrode', 'Elekid', 'Empoleon', 'Feebas', 'Feraligatr', 'Finneon', 'Flaaffy', 'Flareon', 'Floatzel', 
'Flygon', 'Froslass', 'Gabite', 'Garchomp', 'Gastly', 'Gastrodon', 'Gengar', 'Geodude', 'Gible', 'Glaceon', 'Glalie', 
'Gloom', 'Goldeen', 'Golduck', 'Golem', 'Gorebyss', 'Graveler', 'Grimer', 'Grotle', 'Grovyle', 'Growlithe', 'Gulpin', 
'Haunter', 'Hippopotas', 'Hippowdon', 'Horsea', 'Houndoom', 'Houndour', 'Huntail', 'Illumise', 'Infernape', 'Ivysaur', 
'Jolteon', 'Kabuto', 'Kabutops', 'Kakuna', 'Kingdra', 'Kingler', 'Koffing', 'Krabby', 'Kricketot', 'Kricketune', 
'Lanturn', 'Lapras', 'Larvitar', 'Leafeon', 'Lileep', 'Lombre', 'Lotad', 'Ludicolo', 'Lumineon', 'Luvdisc', 'Luxio',
'Luxray', 'Magby', 'Magcargo', 'Magikarp', 'Magmar', 'Magmortar', 'Mamoswine', 'Manectric', 'Mareep', 'Marowak', 
'Marshtomp', 'Meganium', 'Metapod', 'Mightyena', 'Milotic', 'Minun', 'Misdreavus', 'Mismagius', 'Monferno', 'Mudkip',
'Muk', 'Nidoking', 'Nidoqueen', 'Nidoran♀', 'Nidoran♂', 'Nidorina', 'Nidorino', 'Nincada', 'Ninetales', 'Nosepass',
'Numel', 'Nuzleaf', 'Octillery', 'Oddish', 'Omanyte', 'Omastar', 'Onix', 'Pachirisu', 'Paras', 'Parasect', 'Phanpy', 
'Pichu', 'Pikachu', 'Piloswine', 'Pineco', 'Pinsir', 'Piplup', 'Plusle', 'Politoed', 'Poliwag', 'Poliwhirl', 'Ponyta',
'Poochyena', 'Prinplup', 'Psyduck', 'Pupitar', 'Quagsire', 'Quilava', 'Qwilfish', 'Raichu', 'Rampardos', 'Rapidash',
'Relicanth', 'Remoraid', 'Rhydon', 'Rhyhorn', 'Rhyperior', 'Roselia', 'Roserade', 'Sandshrew', 'Sandslash', 'Sceptile',
'Seadra', 'Seaking', 'Sealeo', 'Seedot', 'Seel', 'Seviper', 'Sharpedo', 'Shedinja', 'Shelgon', 'Shellder', 'Shellos',
'Shiftry', 'Shinx', 'Shroomish', 'Shuckle', 'Shuppet', 'Silcoon', 'Skorupi', 'Skuntank', 'Slugma', 'Sneasel', 'Snorunt',
'Snover', 'Spheal', 'Spinarak', 'Spiritomb', 'Squirtle', 'Staryu', 'Stunky', 'Sudowoodo', 'Sunflora', 'Sunkern',
'Surskit', 'Swalot', 'Swampert', 'Swinub', 'Tangela', 'Tangrowth', 'Tentacool', 'Tentacruel', 'Torchic', 'Torkoal',
'Torterra', 'Totodile', 'Trapinch', 'Treecko', 'Turtwig', 'Typhlosion', 'Tyranitar', 'Vaporeon', 'Venomoth', 'Venonat',
'Venusaur', 'Vibrava', 'Victreebel', 'Vileplume', 'Volbeat', 'Voltorb', 'Vulpix', 'Wailmer', 'Wailord', 'Walrein', 
'Wartortle', 'Weavile', 'Weedle', 'Weepinbell', 'Weezing', 'Whiscash', 'Wooper', 'Wormadam (Plant Cloak)', 
'Wormadam (Sandy Cloak)', 'Wurmple') )`)
)

type Cup struct {
	name           string
	pokemon        []dtos.PokemonDto
	moveSets       map[int64]dtos.MoveSetDto
	ids            []int64
	battleMatrix   map[int64]map[int64]float64
	pageRankMatrix *mat.Dense
	mutex          sync.Mutex
	wg             sync.WaitGroup
	current        float64
	goal           float64
}

func (cup *Cup) FillBattleMatrix() {
	cup.battleMatrix = map[int64]map[int64]float64{}
	ids := make(chan int, len(cup.ids))
	for w := 0; w < 10; w++ {
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
		battleSims := daos.BATTLE_SIMS_DAO.FindMatchupsForAlly(ally, cup.ids)
		battleMiniMatrix := map[int64]float64{}
		for _, sim := range battleSims {
			battleMiniMatrix[sim.EnemyId()] = sim.Score()
		}
		cup.mutex.Lock()
		cup.battleMatrix[ally] = battleMiniMatrix
		cup.current += 1.0
		fmt.Printf("%f%% Complete\n", 100.0*cup.current/cup.goal)
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
		err, _ := daos.RANKINGS_DAO.Create(cup.name, ranking.moveSet.PokemonId(), ranking.moveSet.Id(), ranking.pokemonRank, ranking.score)
		daos.CheckError(err)
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

func (cup *Cup) CalculateTeams(pokemonId int64, pokemonRankings map[int64][]Ranking) ([]int64, float64) {
	var bestTeam []int64
	bestScore := 0.0
	pokemon := pokemonRankings[pokemonId][0]
	for _, allyGroupOne := range pokemonRankings {
		allyOne := allyGroupOne[0]
		if allyOne.pokemonRank == pokemon.pokemonRank {
			continue
		}
		teamScore := 0.0
		for _, allyGroupTwo := range pokemonRankings {
			allyTwo := allyGroupTwo[0]
			if allyOne.score >= allyTwo.score {
				continue
			}
			for _, enemyGroup := range pokemonRankings {
				enemy := enemyGroup[0]
				enemyScore := cup.battleMatrix[pokemon.moveSet.Id()][enemy.moveSet.Id()]
				enemyScore = math.Max(enemyScore, cup.battleMatrix[allyOne.moveSet.Id()][enemy.moveSet.Id()])
				enemyScore = math.Max(enemyScore, cup.battleMatrix[allyTwo.moveSet.Id()][enemy.moveSet.Id()])
				teamScore += enemyScore * enemy.score / 100.0
			}
			if bestScore < teamScore {
				bestScore = teamScore
				bestTeam = []int64{allyOne.moveSet.PokemonId(), allyTwo.moveSet.PokemonId()}
			}
		}
	}
	return bestTeam, bestScore
}

func (cup *Cup) CalculateGoodMatchUps(pokemonId int64, pokemonRankings map[int64][]Ranking) []int64 {
	type PokemonAndScore struct {
		pokemonId int64
		score float64
	}
	var enemies []PokemonAndScore

	pokemon := pokemonRankings[pokemonId][0]
	for _, enemyGroup := range pokemonRankings {
		enemy := enemyGroup[0]
		enemyScore := cup.battleMatrix[pokemon.moveSet.Id()][enemy.moveSet.Id()] * enemy.score / 100.0
		enemies = append(enemies, PokemonAndScore{enemy.moveSet.PokemonId(), enemyScore})
	}
	sort.Slice(enemies, func(i, j int) bool {
		return enemies[i].score < enemies[j].score
	})

	var results []int64
	for i := 0; i < 3; i++ {
		results = append(results, enemies[i].pokemonId)
	}
	return results
}

func (cup *Cup) CalculateBadMatchUps(pokemonId int64, pokemonRankings map[int64][]Ranking) []int64 {
	type PokemonAndScore struct {
		pokemonId int64
		score float64
	}
	var enemies []PokemonAndScore

	pokemon := pokemonRankings[pokemonId][0]
	for _, enemyGroup := range pokemonRankings {
		enemy := enemyGroup[0]
		enemyScore := cup.battleMatrix[enemy.moveSet.Id()][pokemon.moveSet.Id()] * enemy.score / 100.0
		enemies = append(enemies, PokemonAndScore{enemy.moveSet.PokemonId(), enemyScore})
	}
	sort.Slice(enemies, func(i, j int) bool {
		return enemies[i].score < enemies[j].score
	})

	var results []int64
	for i := 0; i < 3; i++ {
		results = append(results, enemies[i].pokemonId)
	}
	return results
}

type Ranking struct {
	moveSet     dtos.MoveSetDto
	score       float64
	pokemonRank interface{}
}

func newCup(name, cupSql string) *Cup {
	var cup = Cup{}
	cup.name = name
	cup.pokemon = daos.POKEMON_DAO.FindWhere(cupSql)
	cup.moveSets = map[int64]dtos.MoveSetDto{}
	cup.ids = []int64{}
	for _, pokemon := range cup.pokemon {
		moveSets := daos.MOVE_SETS_DAO.FindWhere("pokemon_id = ? AND simulated", pokemon.Id())
		for _, moveSet := range moveSets {
			cup.moveSets[moveSet.Id()] = moveSet
			cup.ids = append(cup.ids, moveSet.Id())
		}
	}
	cup.mutex = sync.Mutex{}
	cup.wg = sync.WaitGroup{}
	return &cup
}
