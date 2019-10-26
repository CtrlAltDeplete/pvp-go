package db

import (
	"database/sql"
	"testing"
)

func fail(t *testing.T, msg string, expected, actual interface{}) {
	t.Fatalf("%s:\n\tExpected %v\n\tGot %v\n", msg, expected, actual)
}

func fakeTypeSetup() int64 {
	var (
		sqlStmt = "INSERT INTO pvpgo.types (first_type, display_name) " +
			"VALUES ('Test', 'Test')"
		id     int64
		result sql.Result
		err    error
	)
	result, err = LIVE.Exec(sqlStmt)
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	return id
}

func fakeTypeTeardown() {
	var (
		sqlStmt = "DELETE FROM pvpgo.types " +
			"WHERE first_type = 'Test'"
		err error
	)
	_, err = LIVE.Exec(sqlStmt)
	CheckError(err)
}

func fakeTypeMultiplierSetup() (int64, int64) {
	var (
		sqlStmt = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
		typeId int64
		id     int64
		result sql.Result
		err    error
	)
	typeId = fakeTypeSetup()
	result, err = LIVE.Exec(sqlStmt, typeId, typeId, 1.0)
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	return typeId, id
}

func fakeTypeMultiplierTeardown() {
	var (
		sqlStmt = "DELETE tm " +
			"FROM pvpgo.types t " +
			"INNER JOIN pvpgo.type_multipliers tm " +
			"ON t.id = tm.receiving_type OR t.id = tm.acting_type " +
			"WHERE t.first_type = 'Test'"
		err error
	)
	defer fakeTypeTeardown()
	_, err = LIVE.Exec(sqlStmt)
	CheckError(err)
}

func fakePokemonSetup() int64 {
	var (
		sqlStmt = "INSERT INTO pvpgo.pokemon (gen, name, type_id, atk, def, sta) " +
			"VALUES (0, 'Test', 1, 0, 0, 0)"
		id     int64
		result sql.Result
		err    error
	)
	result, err = LIVE.Exec(sqlStmt)
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	return id
}

func fakePokemonTeardown() {
	var (
		sqlStmt = "DELETE FROM pvpgo.pokemon " +
			"WHERE name LIKE 'Test%'"
		err error
	)
	_, err = LIVE.Exec(sqlStmt)
	CheckError(err)
}

func fakeMoveSetup() int64 {
	var (
		sqlStmt = "INSERT INTO pvpgo.moves (name, type_id, power, turns, energy) " +
			"VALUES ('Test', 1, 0, 0, 0)"
		id     int64
		result sql.Result
		err    error
	)
	result, err = LIVE.Exec(sqlStmt)
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	return id
}

func fakeMoveTeardown() {
	var (
		sqlStmt = "DELETE FROM pvpgo.moves " +
			"WHERE name LIKE 'Test%'"
		err error
	)
	_, err = LIVE.Exec(sqlStmt)
	CheckError(err)
}

func fakePokemonHasMoveSetup() (int64, int64, int64) {
	var (
		sqlStmt = "INSERT INTO pvpgo.pokemon_has_move (pokemon_id, move_id) " +
			"VALUES (?, ?)"
		pokemonId int64
		moveId    int64
		id        int64
		result    sql.Result
		err       error
	)
	pokemonId = fakePokemonSetup()
	moveId = fakeMoveSetup()
	result, err = LIVE.Exec(sqlStmt, pokemonId, moveId)
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	return pokemonId, moveId, id
}

func fakePokemonHasMoveTeardown() {
	var (
		pokemonSqlStmt = "DELETE phm " +
			"FROM pvpgo.pokemon p " +
			"INNER JOIN pvpgo.pokemon_has_move phm ON p.id = phm.pokemon_id " +
			"WHERE p.name LIKE 'Test%' "
		moveSqlStmt = "DELETE phm " +
			"FROM pvpgo.moves m " +
			"INNER JOIN pvpgo.pokemon_has_move phm ON m.id = phm.move_id " +
			"WHERE m.name LIKE 'Test%' "
		pokemonErr error
		moveErr    error
	)
	defer fakePokemonTeardown()
	defer fakeMoveTeardown()
	_, pokemonErr = LIVE.Exec(pokemonSqlStmt)
	_, moveErr = LIVE.Exec(moveSqlStmt)
	CheckError(pokemonErr)
	CheckError(moveErr)
}
