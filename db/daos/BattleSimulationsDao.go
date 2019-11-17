package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"time"
)

type BattleSimulationsDao struct{}

func (dao *BattleSimulationsDao) Create(allyId, enemyId int64, individualMatchups []int64) (error, *dtos.BattleSimulationDto) {
	var (
		result sql.Result
		err    error
		id     int64
		score  int64
		query  = "INSERT INTO pvpgo.battle_simulations (ally_id, enemy_id, `0v0`, `0v1`, `0v2`, `1v0`, `1v1`, `1v2`, `2v0`, `2v1`, `2v2`, score) " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	)
	score = dtos.CalculateTotalScore(individualMatchups)
	result, err = LIVE.Exec(query, allyId, enemyId, individualMatchups[0], individualMatchups[1], individualMatchups[2],
		individualMatchups[3], individualMatchups[4], individualMatchups[5], individualMatchups[6],
		individualMatchups[7], individualMatchups[8], score)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newBattleSimulation(id, allyId, enemyId, individualMatchups)
}

func (dao *BattleSimulationsDao) BatchCreate(params []int64) {
	var (
		err                        error
		allyId, enemyId            int64
		zeroZero, zeroOne, zeroTwo int64
		oneZero, oneOne, oneTwo    int64
		twoZero, twoOne, twoTwo    int64
		query                      = "INSERT INTO pvpgo.battle_simulations (ally_id, enemy_id, `0v0`, `0v1`, `0v2`, `1v0`, `1v1`, `1v2`, `2v0`, `2v1`, `2v2`, score) " +
			"VALUES "
	)

	i := 0
	for i < len(params) {
		if i != 0 {
			query += ", \n"
		}
		allyId = params[i]
		i++
		enemyId = params[i]
		i++
		zeroZero = params[i]
		i++
		zeroOne = params[i]
		i++
		zeroTwo = params[i]
		i++
		oneZero = params[i]
		i++
		oneOne = params[i]
		i++
		oneTwo = params[i]
		i++
		twoZero = params[i]
		i++
		twoOne = params[i]
		i++
		twoTwo = params[i]
		i++
		query += fmt.Sprintf("(%d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d, %d)", allyId, enemyId, zeroZero,
			zeroOne, zeroTwo, oneZero, oneOne, oneTwo, twoZero, twoOne, twoTwo, dtos.CalculateTotalScore([]int64{
				zeroZero, zeroOne, zeroTwo, oneZero, oneOne, oneTwo, twoZero, twoOne, twoTwo,
			}))
	}
	_, err = LIVE.Exec(query)
	if err != nil {
		fileName := fmt.Sprintf("%s - error.log", time.Now())
		log.Printf("Error: Creating [%s]\n", fileName)
		file, _ := os.Create(fileName)
		encode := gob.NewEncoder(file)
		_ = encode.Encode(params)
	}
}

func newBattleSimulation(id, allyId, enemyId int64, individualMatchups []int64) *dtos.BattleSimulationDto {
	var battleSim = dtos.BattleSimulationDto{}
	battleSim.SetId(id)
	battleSim.SetAllyId(allyId)
	battleSim.SetEnemyId(enemyId)
	battleSim.SetIndividualMatchups(individualMatchups)
	return &battleSim
}
