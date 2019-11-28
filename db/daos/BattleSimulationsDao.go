package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type BattleSimulationsDao struct{}

func (dao *BattleSimulationsDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.BattleSimulationDto) {
	var (
		id                                          int64
		allyId                                      int64
		enemyId                                     int64
		zvz, zvo, zvt, ovz, ovo, ovt, tvz, tvo, tvt int64
		score                                       int64
		rows                                        *sql.Rows
		err                                         error
		count                                       = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.battle_simulations " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &allyId, &enemyId, &zvz, &zvo, &zvt, &ovz, &ovo, &ovt, &tvz, &tvo, &tvt, &score))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newBattleSimulation(id, allyId, enemyId, []int64{zvz, zvo, zvt, ovz, ovo, ovt, tvz, tvo, tvt})
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *BattleSimulationsDao) FindWhere(query string, params ...interface{}) []dtos.BattleSimulationDto {
	var (
		sims                                        = []dtos.BattleSimulationDto{}
		rows                                        *sql.Rows
		err                                         error
		id                                          int64
		allyId                                      int64
		enemyId                                     int64
		zvz, zvo, zvt, ovz, ovo, ovt, tvz, tvo, tvt int64
		score                                       int64
	)
	query = "SELECT * " +
		"FROM pvpgo.battle_simulations " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &allyId, &enemyId, &zvz, &zvo, &zvt, &ovz, &ovo, &ovt, &tvz, &tvo, &tvt, &score))
		sims = append(sims, *newBattleSimulation(id, allyId, enemyId, []int64{zvz, zvo, zvt, ovz, ovo, ovt, tvz, tvo, tvt}))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return sims
}

func (dao *BattleSimulationsDao) FindMatchupsForAlly(allyId int64, enemyIds []int64) []dtos.BattleSimulationDto {
	var (
		params []interface{}
		query  = "ally_id = ? " +
			"AND enemy_id IN (?" + strings.Repeat(", ?", len(enemyIds)-1) + ")"
	)
	params = append(params, allyId)
	for _, enemyId := range enemyIds {
		params = append(params, enemyId)
	}
	return dao.FindWhere(query, params...)
}

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
	err = errors.New("")
	attempts := 0
	for err != nil && attempts < 10 {
		attempts++
		_, err = LIVE.Exec(query)
	}
	if err != nil {
		fileName := fmt.Sprintf("%s - error.log", time.Now())
		log.Printf("Error [%s]: Creating [%s]\n", err.Error(), fileName)
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
