package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"strings"
)

type MovesDao struct{}

func (dao *MovesDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.MoveDto) {
	var (
		id                           int64
		name                         string
		typeId, power, turns, energy int64
		probability                  sql.NullFloat64
		stageDelta                   sql.NullInt64
		stats, target                sql.NullString
		rows                         *sql.Rows
		err                          error
		count                        = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.moves " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &name, &typeId, &power, &turns, &energy, &probability, &stageDelta, &stats, &target))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newMove(id, name, typeId, power, turns, energy, probability, stageDelta, stats, target)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *MovesDao) FindById(id int64) (error, *dtos.MoveDto) {
	var (
		query = "id = ?"
	)
	return dao.FindSingleWhere(query, id)
}

func (dao *MovesDao) FindByName(name string) (error, *dtos.MoveDto) {
	var (
		query = "name = ?"
	)
	return dao.FindSingleWhere(query, name)
}

func (dao *MovesDao) FindWhere(query string, params ...interface{}) []dtos.MoveDto {
	var (
		moves                        = []dtos.MoveDto{}
		rows                         *sql.Rows
		e                            error
		id                           int64
		name                         string
		typeId, power, turns, energy int64
		probability                  sql.NullFloat64
		stageDelta                   sql.NullInt64
		stats, target                sql.NullString
	)
	query = "SELECT * " +
		"FROM pvpgo.moves " +
		"WHERE " + query
	rows, e = LIVE.Query(query, params...)
	CheckError(e)
	for rows.Next() {
		CheckError(rows.Scan(&id, &name, &typeId, &power, &turns, &energy, &probability, &stageDelta, &stats, &target))
		moves = append(moves, *newMove(id, name, typeId, power, turns, energy, probability, stageDelta, stats, target))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return moves
}

func (dao *MovesDao) FindByTypeId(id int64) []dtos.MoveDto {
	var (
		query = "type_id = ?"
	)
	return dao.FindWhere(query, id)
}

func (dao *MovesDao) FindByTypeIds(ids []int64) []dtos.MoveDto {
	var (
		id     int64
		params []interface{}
		query  = "type_id IN (?" + strings.Repeat(", ?", len(ids)-1) + ")"
	)
	for _, id = range ids {
		params = append(params, id)
	}
	return dao.FindWhere(query, params...)
}

func (dao *MovesDao) FindAll() []dtos.MoveDto {
	return dao.FindWhere("TRUE")
}

func (dao *MovesDao) Create(name string, typeId, power, turns, energy int64, probability, stageDelta, stats,
	target interface{}) (error, *dtos.MoveDto) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.moves (name, type_id, power, turns, energy, probability, stage_delta, stats, target) " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	)
	result, err = LIVE.Exec(query, name, typeId, power, turns, energy, probability, stageDelta, stats, target)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newMove(id, name, typeId, power, turns, energy, probability, stageDelta, stats, target)
}

func (dao *MovesDao) Update(move dtos.MoveDto) {
	var (
		e     error
		query = "UPDATE pvpgo.moves " +
			"SET name = ?, " +
			"type_id = ?, " +
			"power = ?, " +
			"turns = ?, " +
			"energy = ?, " +
			"probability = ?, " +
			"stage_delta = ?, " +
			"stats = ?, " +
			"target = ? " +
			"WHERE id = ?"
	)
	_, e = LIVE.Exec(query, move.Name(), move.TypeId(), int64(move.Power()), move.Turns(), move.Energy(),
		move.ProbabilityNullable(), move.StageDeltaNullable(), move.StatsNullable(), move.TargetNullable(), move.Id())
	CheckError(e)
}

func (dao *MovesDao) Upsert(name string, typeId, power, turns, energy int64, probability, stageDelta, stats,
	target interface{}) (error, *dtos.MoveDto) {
	var (
		err  error
		move *dtos.MoveDto
	)
	err, move = dao.FindByName(name)
	if err == NO_ROWS {
		err, move = dao.Create(name, typeId, power, turns, energy, probability, stageDelta, stats, target)
	} else if err == nil {
		move.SetTypeId(typeId)
		move.SetPower(float64(power))
		move.SetTurns(turns)
		move.SetEnergy(energy)
		move.SetProbability(probability)
		move.SetStageDelta(stageDelta)
		move.SetStats(stats)
		move.SetTarget(target)
		dao.Update(*move)
	}
	if err != nil {
		return err, nil
	}
	return nil, move
}

func (dao *MovesDao) Delete(move dtos.MoveDto) {
	var (
		e     error
		query = "DELETE FROM pvpgo.moves " +
			"WHERE id = ?"
	)
	_, e = LIVE.Exec(query, move.Id())
	CheckError(e)
}

func newMove(id int64, name string, typeId int64, power int64, turns int64, energy int64, probability interface{},
	stageDelta interface{}, stats interface{}, target interface{}) *dtos.MoveDto {
	var m = dtos.MoveDto{}
	m.SetId(id)
	m.SetName(name)
	m.SetTypeId(typeId)
	m.SetPower(float64(power))
	m.SetTurns(turns)
	m.SetEnergy(energy)
	m.SetProbability(probability)
	m.SetStageDelta(stageDelta)
	m.SetStats(stats)
	m.SetTarget(target)
	return &m
}
