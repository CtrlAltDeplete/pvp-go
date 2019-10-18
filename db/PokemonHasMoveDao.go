package db

import (
	"PvP-Go/models"
	"database/sql"
)

type PokemonHasMoveDao struct{}

func (dao *PokemonHasMoveDao) FindByPokemonAndMove(pokemonId, moveId int64) (error, *models.PokemonHasMove) {
	var (
		id       int64
		isLegacy bool
		query    = "SELECT * " +
			"FROM pvpgo.pokemon_has_move " +
			"WHERE pokemon_id = ? " +
			"AND move_id = ?"
		err error
	)
	err = LIVE.QueryRow(query, pokemonId, moveId).Scan(&id, &pokemonId, &moveId, &isLegacy)
	if err != nil {
		return err, nil
	}
	return nil, newPokemonHasMove(id, pokemonId, moveId, isLegacy)
}

func (dao *PokemonHasMoveDao) Create(pokemonId, moveId int64, isLegacy bool) *models.PokemonHasMove {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.pokemon_has_move (pokemon_id, move_id, is_legacy) " +
			"VALUES (?, ?, ?)"
	)
	result, err = LIVE.Exec(query, pokemonId, moveId, isLegacy)
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	return newPokemonHasMove(id, pokemonId, moveId, isLegacy)
}

func (dao *PokemonHasMoveDao) FindOrCreate(pokemonId, moveId int64, isLegacy bool) *models.PokemonHasMove {
	var (
		pokemonHasMove *models.PokemonHasMove
		err            error
	)
	err, pokemonHasMove = dao.FindByPokemonAndMove(pokemonId, moveId)
	if err != nil {
		return dao.Create(pokemonId, moveId, isLegacy)
	}
	return pokemonHasMove
}

func newPokemonHasMove(id, pokemonId, moveId int64, isLegacy bool) *models.PokemonHasMove {
	var phm = models.PokemonHasMove{}
	phm.SetId(id)
	phm.SetPokemonId(pokemonId)
	phm.SetMoveId(moveId)
	phm.SetIsLegacy(isLegacy)
	return &phm
}
