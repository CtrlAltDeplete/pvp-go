package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
)

type MoveSetDao struct{}

func (dao *MoveSetDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.MoveSetDto) {
	var (
		id                    int64
		pokemonId             int64
		fastMoveId            int64
		primaryChargeMoveId   int64
		secondaryChargeMoveId *sql.NullInt64
		rows                  *sql.Rows
		err                   error
		count                 = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.move_sets " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &pokemonId, &fastMoveId, &primaryChargeMoveId, &secondaryChargeMoveId))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newMoveSet(id, pokemonId, fastMoveId, primaryChargeMoveId, secondaryChargeMoveId)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *MoveSetDao) FindWhere(query string, params ...interface{}) []dtos.MoveSetDto {
	var (
		moveSets              = []dtos.MoveSetDto{}
		rows                  *sql.Rows
		err                   error
		id                    int64
		pokemonId             int64
		fastMoveId            int64
		primaryChargeMoveId   int64
		secondaryChargeMoveId *sql.NullInt64
		simulated             bool
	)
	query = "SELECT * " +
		"FROM pvpgo.move_sets " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &pokemonId, &fastMoveId, &primaryChargeMoveId, &secondaryChargeMoveId, &simulated))
		moveSets = append(moveSets, *newMoveSet(id, pokemonId, fastMoveId, primaryChargeMoveId, secondaryChargeMoveId, simulated))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return moveSets
}

func (dao *MoveSetDao) FindAll() []dtos.MoveSetDto {
	return dao.FindWhere("TRUE ORDER BY id ASC")
}

func (dao *MoveSetDao) Update(moveSet dtos.MoveSetDto) {
	var (
		err   error
		query = "UPDATE pvpgo.move_sets " +
			"SET simulated = ? " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, moveSet.Simulated(), moveSet.Id())
	CheckError(err)
}

func newMoveSet(id, pokemonId, fastMoveId, primaryChargeMoveId int64, secondaryChargeMoveId *sql.NullInt64, simulated bool) *dtos.MoveSetDto {
	var moveSet = dtos.MoveSetDto{}
	moveSet.SetId(id)
	moveSet.SetPokemonId(pokemonId)
	moveSet.SetFastMoveId(fastMoveId)
	moveSet.SetPrimaryChargeMoveId(primaryChargeMoveId)
	if secondaryChargeMoveId != nil && secondaryChargeMoveId.Valid {
		moveSet.SetSecondaryChargeMoveId(&secondaryChargeMoveId.Int64)
	} else {
		moveSet.SetSecondaryChargeMoveId(nil)
	}
	moveSet.SetSimulated(simulated)
	return &moveSet
}
