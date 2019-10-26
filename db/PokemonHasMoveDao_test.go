package db

import (
	"PvP-Go/models"
	"database/sql"
	"reflect"
	"testing"
)

func TestPokemonHasMoveDao_FindSingleWhere(t *testing.T) {
	// Initialize test variables
	var (
		query            string
		moveId           int64
		pokemonId        int64
		pokemonHasMoveId int64
		result           sql.Result
		err              error
		expected         *models.PokemonHasMove
		actual           *models.PokemonHasMove
		setupSql         = "INSERT INTO pvpgo.pokemon_has_move (pokemon_id, move_id) " +
			"VALUES (?, ?)"
	)

	// Defer teardown
	defer fakePokemonHasMoveTeardown()

	// Test setup
	pokemonId = fakePokemonSetup()
	moveId = fakeMoveSetup()
	result, err = LIVE.Exec(setupSql, pokemonId, moveId)
	CheckError(err)
	pokemonHasMoveId, err = result.LastInsertId()
	CheckError(err)
	_, err = LIVE.Exec(setupSql, pokemonId, 1)
	CheckError(err)

	/* Happy Path Test */
	// Prepare test variables
	query = "id <= ? " +
		"AND id >= ?"
	expected = newPokemonHasMove(pokemonHasMoveId, pokemonId, moveId, false)
	err, actual = POKEMON_HAS_MOVE_DAO.FindSingleWhere(query, pokemonHasMoveId, pokemonHasMoveId)
	CheckError(err)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "PokemonHasMoveDao.FindSingleWhere failed during the happy path", *expected, *actual)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	err, _ = POKEMON_HAS_MOVE_DAO.FindSingleWhere(query, 0)

	// Check expected vs actual
	if err != NO_ROWS {
		fail(t, "PokemonHasMoveDao.FindSingleWhere failed on no results", NO_ROWS, err)
	}

	/* Multiple Results Test */
	// Prepare test variables
	query = "id >= ?"
	err, _ = POKEMON_HAS_MOVE_DAO.FindSingleWhere(query, pokemonHasMoveId)

	// Check expected vs actual
	if err != MULTIPLE_ROWS {
		fail(t, "PokemonHasMoveDao.FindSingleWhere failed on multiple results", MULTIPLE_ROWS, err)
	}
}

func TestPokemonHasMoveDao_FindWhere(t *testing.T) {
	// Initialize test variables
	var (
		query             string
		pokemonId         int64
		id                int64
		pokemonHasMoveIds []int64
		result            sql.Result
		err               error
		expected          []models.PokemonHasMove
		actual            []models.PokemonHasMove
		setupSql          = "INSERT INTO pvpgo.pokemon_has_move (pokemon_id, move_id) " +
			"VALUES (?, ?)"
	)

	// Defer teardown
	defer fakePokemonHasMoveTeardown()

	// Test setup
	pokemonId = fakePokemonSetup()
	for i := int64(1); i < 3; i++ {
		result, err = LIVE.Exec(setupSql, pokemonId, i)
		CheckError(err)
		id, err = result.LastInsertId()
		CheckError(err)
		pokemonHasMoveIds = append(pokemonHasMoveIds, id)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	expected = []models.PokemonHasMove{}
	actual = POKEMON_HAS_MOVE_DAO.FindWhere(query, 0)

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "PokemonHasMoveDao.FindWhere failed on no results", expected, actual)
	}

	/* Multiple Results Test */
	// Prepare test variables
	if pokemonHasMoveIds == nil {
		panic("Expected non-nil pokemonHasMoveIds")
	}
	query = "id IN (?, ?)"
	expected = []models.PokemonHasMove{
		*newPokemonHasMove(pokemonHasMoveIds[0], pokemonId, 1, false),
		*newPokemonHasMove(pokemonHasMoveIds[1], pokemonId, 2, false),
	}
	actual = POKEMON_HAS_MOVE_DAO.FindWhere(query, pokemonHasMoveIds[0], pokemonHasMoveIds[1])

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "PokemonHasMoveDao.FindWhere failed on multiple results", expected, actual)
	}
}

func TestPokemonHasMoveDao_Create(t *testing.T) {
	// Initialize test variables
	var (
		pokemonId int64
		moveId    int64
		err       error
		expected  *models.PokemonHasMove
		actual    *models.PokemonHasMove
	)

	// Defer teardown
	defer fakePokemonHasMoveTeardown()

	// Prepare test variables
	pokemonId = fakePokemonSetup()
	moveId = fakeMoveSetup()
	err, actual = POKEMON_HAS_MOVE_DAO.Create(pokemonId, moveId, false)
	CheckError(err)
	expected = newPokemonHasMove(actual.Id(), pokemonId, moveId, false)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "PokemonHasMoveDao.Create failed", *expected, *actual)
	}
}

func TestPokemonHasMoveDao_Update(t *testing.T) {
	// Initialize test variables
	var (
		pokemonId        int64
		moveId           int64
		pokemonHasMoveId int64
		isLegacy         bool
		err              error
		expected         *models.PokemonHasMove
		actual           *models.PokemonHasMove
		verifySql        = "SELECT * " +
			"FROM pvpgo.pokemon_has_move " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakePokemonHasMoveTeardown()

	// Prepare test variables
	pokemonId, moveId, pokemonHasMoveId = fakePokemonHasMoveSetup()
	expected = newPokemonHasMove(pokemonHasMoveId, pokemonId, moveId, true)
	POKEMON_HAS_MOVE_DAO.Update(*expected)
	err = LIVE.QueryRow(verifySql, pokemonHasMoveId).Scan(&pokemonHasMoveId, &pokemonId, &moveId, &isLegacy)
	CheckError(err)
	actual = newPokemonHasMove(pokemonHasMoveId, pokemonId, moveId, isLegacy)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "PokemonHasMoveDao.Update failed", *expected, *actual)
	}
}

func TestPokemonHasMoveDao_Delete(t *testing.T) {
	// Initialize test variables
	var (
		pokemonId        int64
		moveId           int64
		pokemonHasMoveId int64
		pokemonHasMove   *models.PokemonHasMove
		expected         int64 = 0
		actual           int64
		verifySql        = "SELECT COUNT(*) " +
			"FROM pvpgo.pokemon_has_move " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakePokemonHasMoveTeardown()

	// Prepare test variables
	pokemonId, moveId, pokemonHasMoveId = fakePokemonHasMoveSetup()
	pokemonHasMove = newPokemonHasMove(pokemonHasMoveId, pokemonId, moveId, false)
	POKEMON_HAS_MOVE_DAO.Delete(*pokemonHasMove)
	CheckError(LIVE.QueryRow(verifySql, pokemonHasMoveId).Scan(&actual))

	// Check expected vs actual
	if expected != actual {
		fail(t, "PokemonHasMoveDao.Delete failed", expected, actual)
	}
}
