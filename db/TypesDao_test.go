package db

import (
	"PvP-Go/models"
	"database/sql"
	"reflect"
	"testing"
)

func TestTypesDao_FindSingleWhere(t *testing.T) {
	// Initialize test variables
	var (
		query     string
		typeId    int64
		typeNames = []string{"Test1", "Test2"}
		result    sql.Result
		err       error
		expected  *models.PokemonType
		actual    *models.PokemonType
		setupSql  = "INSERT INTO pvpgo.types (first_type, display_name) " +
			"VALUES (?, ?)"
		teardownSql = "DELETE FROM pvpgo.types " +
			"WHERE display_name IN (?, ?)"
	)

	// Defer teardown
	defer func() {
		_, err = LIVE.Exec(teardownSql, typeNames[0], typeNames[1])
		CheckError(err)
	}()

	// Test setup
	result, err = LIVE.Exec(setupSql, typeNames[0], typeNames[0])
	CheckError(err)
	typeId, err = result.LastInsertId()
	CheckError(err)
	_, err = LIVE.Exec(setupSql, typeNames[1], typeNames[1])
	CheckError(err)

	/* Happy Path Test */
	// Prepare test variables
	query = "id <= ? " +
		"AND id >= ?"
	expected = newPokemonType(typeId, typeNames[0], nil, typeNames[0])
	err, actual = TYPES_DAO.FindSingleWhere(query, typeId, typeId)
	CheckError(err)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypesDao.FindSingleWhere failed during the happy path", *expected, *actual)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	err, _ = TYPES_DAO.FindSingleWhere(query, 0)

	// Check expected vs actual
	if err != NO_ROWS {
		fail(t, "TypesDao.FindSingleWhere failed on no results", NO_ROWS, err)
	}

	/* Multiplier Results Test */
	// Prepare test variables
	query = "id >= ?"
	err, _ = TYPES_DAO.FindSingleWhere(query, typeId)

	// Check expected vs actual
	if err != MULTIPLE_ROWS {
		fail(t, "TypesDao.FindSingleWhere failed on multiple results", MULTIPLE_ROWS, err)
	}
}

func TestTypesDao_FindWhere(t *testing.T) {
	// Initialize test variables
	var (
		query     string
		id        int64
		typeIds   []int64
		typeNames = []string{"Test1", "Test2"}
		result    sql.Result
		err       error
		expected  []models.PokemonType
		actual    []models.PokemonType
		setupSql  = "INSERT INTO pvpgo.types (first_type, display_name) " +
			"VALUES (?, ?)"
		teardownSql = "DELETE FROM pvpgo.types " +
			"WHERE display_name IN (?, ?)"
	)

	// Defer teardown
	defer func() {
		_, err = LIVE.Exec(teardownSql, typeNames[0], typeNames[1])
		CheckError(err)
	}()

	// Test setup
	result, err = LIVE.Exec(setupSql, typeNames[0], typeNames[0])
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	typeIds = append(typeIds, id)
	result, err = LIVE.Exec(setupSql, typeNames[1], typeNames[1])
	CheckError(err)
	id, err = result.LastInsertId()
	CheckError(err)
	typeIds = append(typeIds, id)

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	expected = []models.PokemonType{}
	actual = TYPES_DAO.FindWhere(query, 0)

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "TypesDao.FindWhere failed on no results", expected, actual)
	}

	/* Multiple Results Test */
	// Prepare test variables
	query = "id IN (?, ?)"
	expected = []models.PokemonType{
		*newPokemonType(typeIds[0], typeNames[0], nil, typeNames[0]),
		*newPokemonType(typeIds[1], typeNames[1], nil, typeNames[1]),
	}
	actual = TYPES_DAO.FindWhere(query, typeIds[0], typeIds[1])

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "TypesDao.FindWhere failed on multiple results", expected, actual)
	}
}

func TestTypesDao_Create(t *testing.T) {
	// Initialize test variables
	var (
		firstType  = "Test"
		secondType sql.NullString
		err        error
		expected   *models.PokemonType
		actual     *models.PokemonType
	)

	// Defer teardown
	defer fakeTypeTeardown()

	/* string/string */
	// Prepare test variables
	err, actual = TYPES_DAO.Create(firstType, firstType)
	CheckError(err)
	expected = newPokemonType(actual.Id(), firstType, firstType, firstType+"/"+firstType)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypesDao.Create failed on string/string", *expected, *actual)
	}

	fakeTypeTeardown()

	/* string/nil */
	// Prepare test variables
	err, actual = TYPES_DAO.Create(firstType, nil)
	CheckError(err)
	expected = newPokemonType(actual.Id(), firstType, nil, firstType)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypesDao.Create failed on string/nil", *expected, *actual)
	}

	fakeTypeTeardown()

	/* string/sql.NullString(valid) */
	// Prepare test variables
	secondType = sql.NullString{
		String: firstType,
		Valid:  true,
	}
	err, actual = TYPES_DAO.Create(firstType, secondType)
	CheckError(err)
	expected = newPokemonType(actual.Id(), firstType, firstType, firstType+"/"+firstType)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypesDao.Create failed on sql.NullString(valid)", *expected, *actual)
	}

	fakeTypeTeardown()

	/* string/sql.NullString(invalid) */
	// Prepare test variables
	secondType = sql.NullString{
		String: "",
		Valid:  false,
	}
	err, actual = TYPES_DAO.Create(firstType, secondType)
	CheckError(err)
	expected = newPokemonType(actual.Id(), firstType, nil, firstType)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypesDao.Create failed on sql.NullString(invalid)", *expected, *actual)
	}

	fakeTypeTeardown()

	/* string/something */
	// Prepare test variables
	err, actual = TYPES_DAO.Create(firstType, float64(0))

	// Check expected vs actual
	if err != BAD_PARAMS {
		fail(t, "TypesDao.Create failed on BAD_PARAMS", BAD_PARAMS, err)
	}
}

func TestTypesDao_Update(t *testing.T) {
	// Initialize test variables
	var (
		firstType   = "Test"
		secondType  sql.NullString
		displayName string
		typeId      int64
		err         error
		expected    *models.PokemonType
		actual      *models.PokemonType
		verifySql   = "SELECT id, first_type, second_type, display_name " +
			"FROM pvpgo.types " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakeTypeTeardown()

	// Prepare test variables
	typeId = fakeTypeSetup()
	expected = newPokemonType(typeId, firstType, firstType, firstType+"/"+firstType)
	TYPES_DAO.Update(*expected)
	err = LIVE.QueryRow(verifySql, typeId).Scan(&typeId, &firstType, &secondType, &displayName)
	CheckError(err)
	actual = newPokemonType(typeId, firstType, secondType, displayName)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypesDao.Update failed", *expected, *actual)
	}
}

func TestTypesDao_Delete(t *testing.T) {
	// Initialize test variables
	var (
		firstType   = "Test"
		typeId      int64
		pokemonType *models.PokemonType
		expected    int64 = 0
		actual      int64
		verifySql   = "SELECT COUNT(*) " +
			"FROM pvpgo.types " +
			"WHERE first_type = ?"
	)

	// Defer teardown
	defer fakeTypeTeardown()

	// Prepare test variables
	typeId = fakeTypeSetup()
	pokemonType = newPokemonType(typeId, firstType, nil, firstType)
	TYPES_DAO.Delete(*pokemonType)
	CheckError(LIVE.QueryRow(verifySql, firstType).Scan(&actual))

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "TypesDao.Delete failed", expected, actual)
	}
}
