package db

import (
	"PvP-Go/models"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func TestPokemonDao_FindSingleWhere(t *testing.T) {
	// Initialize test variables
	var (
		query     string
		typeId    int64
		pokemonId int64
		result    sql.Result
		err       error
		expected  *models.Pokemon
		actual    *models.Pokemon
		setupSql  = "INSERT INTO pvpgo.pokemon (gen, name, type_id, atk, def, sta) " +
			"VALUES (0, ?, ?, 0, 0, 0)"
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakePokemonTeardown()

	// Test setup
	typeId = fakeTypeSetup()
	result, err = LIVE.Exec(setupSql, "Test1", typeId)
	CheckError(err)
	pokemonId, err = result.LastInsertId()
	CheckError(err)
	_, err = LIVE.Exec(setupSql, "Test2", typeId)
	CheckError(err)

	/* Happy Path Test */
	// Prepare test variables
	query = "id <= ? " +
		"AND id >= ?"
	expected = newPokemon(pokemonId, 0, "Test1", typeId, 0, 0, 0, "2017-07-06",
		false, true, 0, 0, 0, 0)
	err, actual = POKEMON_DAO.FindSingleWhere(query, pokemonId, pokemonId)
	CheckError(err)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "PokemonDao.FindSingleWhere failed during the happy path", *expected, *actual)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	err, _ = POKEMON_DAO.FindSingleWhere(query, 0)

	// Check expected vs actual
	if err != NO_ROWS {
		fail(t, "PokemonDao.FindSingleWhere failed on no results", NO_ROWS, err)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id >= ?"
	err, _ = POKEMON_DAO.FindSingleWhere(query, pokemonId)

	// Check expected vs actual
	if err != MULTIPLE_ROWS {
		fail(t, "PokemonDao.FindSingleWhere failed on multiple results", MULTIPLE_ROWS, err)
	}
}

func TestPokemonDao_FindWhere(t *testing.T) {
	// Initialize test variables
	var (
		query      string
		typeId     int64
		id         int64
		pokemonIds []int64
		result     sql.Result
		err        error
		expected   []models.Pokemon
		actual     []models.Pokemon
		setupSql   = "INSERT INTO pvpgo.pokemon (gen, name, type_id, atk, def, sta) " +
			"VALUES (0, ?, ?, 0, 0, 0)"
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakePokemonTeardown()

	// Test setup
	typeId = fakeTypeSetup()
	for i := int64(1); i < 3; i++ {
		result, err = LIVE.Exec(setupSql, fmt.Sprintf("Test%d", i), typeId)
		CheckError(err)
		id, err = result.LastInsertId()
		CheckError(err)
		pokemonIds = append(pokemonIds, id)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	expected = []models.Pokemon{}
	actual = POKEMON_DAO.FindWhere(query, 0)

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "PokemonDao.FindWhere failed on no results", expected, actual)
	}

	/* Multiple Results Test */
	// Prepare test variables
	if pokemonIds == nil {
		panic("Expected non-nil pokemonIds")
	}
	query = "id IN (?, ?)"
	expected = []models.Pokemon{
		*newPokemon(pokemonIds[0], 0, "Test1", typeId, 0, 0, 0, "2017-07-06",
			false, true, 0, 0, 0, 0),
		*newPokemon(pokemonIds[1], 0, "Test2", typeId, 0, 0, 0, "2017-07-06",
			false, true, 0, 0, 0, 0),
	}
	actual = POKEMON_DAO.FindWhere(query, pokemonIds[0], pokemonIds[1])

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "PokemonDao.FindWhere failed on multiple results", expected, actual)
	}
}

func TestPokemonDao_Create(t *testing.T) {
	// Initialize test variables
	var (
		typeId   int64
		err      error
		expected *models.Pokemon
		actual   *models.Pokemon
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakePokemonTeardown()

	// Prepare test variables
	typeId = fakeTypeSetup()
	err, actual = POKEMON_DAO.Create(0, "Test", typeId, 0, 0, 0, "2017-07-06",
		false, true, 0, 0, 0, 0)
	CheckError(err)
	expected = newPokemon(actual.Id(), 0, "Test", typeId, 0, 0, 0, "2017-07-06",
		false, true, 0, 0, 0, 0)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "PokemonDao.Create failed", *expected, *actual)
	}
}

func TestPokemonDao_Update(t *testing.T) {
	// Initialize test variables
	var (
		typeId      int64
		pokemonId   int64
		gen         int64
		name        string
		atk         float64
		def         float64
		sta         float64
		dateAdd     string
		legendary   bool
		pvpEligible bool
		optLevel    float64
		optAtk      float64
		optDef      float64
		optSta      float64
		err         error
		expected    *models.Pokemon
		actual      *models.Pokemon
		verifySql   = "SELECT * " +
			"FROM pvpgo.pokemon " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakePokemonTeardown()

	// Prepare test variables
	typeId = fakeTypeSetup()
	pokemonId = fakePokemonSetup()
	expected = newPokemon(pokemonId, 0, "Test", typeId, 0, 0, 0, "2017-07-06",
		false, true, 0, 0, 0, 0)
	POKEMON_DAO.Update(*expected)
	err = LIVE.QueryRow(verifySql, pokemonId).Scan(&pokemonId, &gen, &name, &typeId, &atk, &def, &sta, &dateAdd,
		&legendary, &pvpEligible, &optLevel, &optAtk, &optDef, &optSta)
	CheckError(err)
	actual = newPokemon(pokemonId, gen, name, typeId, atk, def, sta, dateAdd,
		legendary, pvpEligible, optLevel, optAtk, optDef, optSta)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "PokemonDao.Update failed", *expected, *actual)
	}
}

func TestPokemonDao_Delete(t *testing.T) {
	// Initialize test variables
	var (
		pokemonId int64
		pokemon   *models.Pokemon
		expected  int64 = 0
		actual    int64
		verifySql = "SELECT COUNT(*) " +
			"FROM pvpgo.pokemon " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakePokemonTeardown()

	// Prepare test variables
	pokemonId = fakePokemonSetup()
	pokemon = newPokemon(pokemonId, 0, "Test", 0, 0, 0, 0, "2017-07-06",
		false, true, 0, 0, 0, 0)
	POKEMON_DAO.Delete(*pokemon)
	CheckError(LIVE.QueryRow(verifySql, pokemonId).Scan(&actual))

	// Check expected vs actual
	if expected != actual {
		fail(t, "PokemonDao.Delete failed", expected, actual)
	}
}
