package db

import (
	"PvP-Go/models"
	"database/sql"
)

type PokemonHasMoveDao struct{}

func (dao *PokemonHasMoveDao) FindSingleWhere(query string, params ...interface{}) (error, *models.PokemonHasMove) {
	var (
		id        int64
		pokemonId int64
		moveId    int64
		legacy    bool
		rows      *sql.Rows
		err       error
		count     = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.pokemon_has_move " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &pokemonId, &moveId, &legacy))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newPokemonHasMove(id, pokemonId, moveId, legacy)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *PokemonHasMoveDao) FindByPokemonAndMove(pokemonId, moveId int64) (error, *models.PokemonHasMove) {
	var (
		query = "pokemon_id = ? " +
			"AND move_id = ?"
	)
	return dao.FindSingleWhere(query, pokemonId, moveId)
}

func (dao *PokemonHasMoveDao) FindWhere(query string, params ...interface{}) []models.PokemonHasMove {
	var (
		pokemonHasMoves = []models.PokemonHasMove{}
		rows            *sql.Rows
		err             error
		id              int64
		pokemonId       int64
		moveId          int64
		legacy          bool
	)
	query = "SELECT * " +
		"FROM pvpgo.pokemon_has_move " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &pokemonId, &moveId, &legacy))
		pokemonHasMoves = append(pokemonHasMoves, *newPokemonHasMove(id, pokemonId, moveId, legacy))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return pokemonHasMoves
}

func (dao *PokemonHasMoveDao) FindAllByPokemonId(pokemonId int64) []models.PokemonHasMove {
	var (
		query = "pokemon_id = ?"
	)
	return dao.FindWhere(query, pokemonId)
}

func (dao *PokemonHasMoveDao) Create(pokemonId, moveId int64, isLegacy bool) (error, *models.PokemonHasMove) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.pokemon_has_move (pokemon_id, move_id, is_legacy) " +
			"VALUES (?, ?, ?)"
	)
	result, err = LIVE.Exec(query, pokemonId, moveId, isLegacy)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newPokemonHasMove(id, pokemonId, moveId, isLegacy)
}

func (dao *PokemonHasMoveDao) Update(pokemonHasMove models.PokemonHasMove) {
	var (
		err   error
		query = "UPDATE pvpgo.pokemon_has_move " +
			"SET pokemon_id = ?, " +
			"move_id = ?, " +
			"is_legacy = ? " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, pokemonHasMove.PokemonId(), pokemonHasMove.MoveId(), pokemonHasMove.IsLegacy(), pokemonHasMove.Id())
	CheckError(err)
}

func (dao *PokemonHasMoveDao) Upsert(pokemonId, moveId int64, isLegacy bool) (error, *models.PokemonHasMove) {
	var (
		err            error
		pokemonHasMove *models.PokemonHasMove
	)
	err, pokemonHasMove = dao.FindByPokemonAndMove(pokemonId, moveId)
	if err == NO_ROWS {
		err, pokemonHasMove = dao.Create(pokemonId, moveId, isLegacy)
	} else if err == nil {
		pokemonHasMove.SetIsLegacy(isLegacy)
		dao.Update(*pokemonHasMove)
	}
	if err != nil {
		return err, nil
	}
	return nil, pokemonHasMove
}

func (dao *PokemonHasMoveDao) Delete(pokemonHasMove models.PokemonHasMove) {
	var (
		err   error
		query = "DELETE FROM pvpgo.pokemon_has_move " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, pokemonHasMove.Id())
	CheckError(err)
}

func newPokemonHasMove(id, pokemonId, moveId int64, isLegacy bool) *models.PokemonHasMove {
	var phm = models.PokemonHasMove{}
	phm.SetId(id)
	phm.SetPokemonId(pokemonId)
	phm.SetMoveId(moveId)
	phm.SetIsLegacy(isLegacy)
	return &phm
}
