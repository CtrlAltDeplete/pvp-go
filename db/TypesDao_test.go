package db

import (
	"PvP-Go/models"
	"database/sql"
	"reflect"
	"testing"
)

type typesDaoSingleTestCase struct {
	query    string
	params   []interface{}
	err      error
	expected *models.PokemonType
	actual   *models.PokemonType
}

func TestTypesDao_FindSingleWhere(t *testing.T) {
	testTypesDao_FindSingleWhere_HappyPath(t)
	testTypesDao_FindSingleWhere_NoResults(t)
	testTypesDao_FindSingleWhere_MultipleResults(t)
}

func testTypesDao_FindSingleWhere_HappyPath(t *testing.T) {
	// Initialize test case's variables
	var testCase = typesDaoSingleTestCase{
		query:    "id > ? AND id < ?",
		params:   nil,
		err:      nil,
		expected: newPokemonType(4, "water", nil, "water"),
		actual:   nil,
	}
	testCase.params = append(testCase.params, 3, 5)

	// Get result
	testCase.err, testCase.actual = TYPES_DAO.FindSingleWhere(testCase.query, testCase.params...)
	CheckError(testCase.err)

	// Check expected vs actual
	if *testCase.expected != *testCase.actual {
		fail(t, "TypesDao.FindSingleWhere", *testCase.expected, *testCase.actual)
	}
}

func testTypesDao_FindSingleWhere_NoResults(t *testing.T) {
	// Initialize test case's variables
	var testCase = typesDaoSingleTestCase{
		query:    "id < ?",
		params:   nil,
		err:      nil,
		expected: nil,
		actual:   nil,
	}
	testCase.params = append(testCase.params, 1)

	// Get result
	testCase.err, testCase.actual = TYPES_DAO.FindSingleWhere(testCase.query, testCase.params...)

	// Expected an error
	if testCase.err != NO_ROWS {
		fail(t, "TypesDao.FindSingleWhere", NO_ROWS, testCase.err)
	}
}

func testTypesDao_FindSingleWhere_MultipleResults(t *testing.T) {
	// Initialize test case's variables
	var testCase = typesDaoSingleTestCase{
		query:    "id < ?",
		params:   nil,
		err:      nil,
		expected: nil,
		actual:   nil,
	}
	testCase.params = append(testCase.params, 4)

	// Get result
	testCase.err, testCase.actual = TYPES_DAO.FindSingleWhere(testCase.query, testCase.params...)

	// Expected an error
	if testCase.err != MULTIPLE_ROWS {
		fail(t, "TypesDao.FindSingleWhere", MULTIPLE_ROWS, testCase.err)
	}
}

type typesDaoMultipleTestCase struct {
	query    string
	params   []interface{}
	expected []models.PokemonType
	actual   []models.PokemonType
}

func TestTypesDao_FindWhere(t *testing.T) {
	testTypesDao_FindWhere_NoResults(t)
	testTypesDao_FindWhere_MultipleResults(t)
}

func testTypesDao_FindWhere_NoResults(t *testing.T) {
	// Initialize test case's variables
	var testCase = typesDaoMultipleTestCase{
		query:    "id < ?",
		params:   nil,
		expected: []models.PokemonType{},
		actual:   nil,
	}
	testCase.params = append(testCase.params, 1)

	// Get result
	testCase.actual = TYPES_DAO.FindWhere(testCase.query, testCase.params...)

	// Expected vs Actual
	if !reflect.DeepEqual(testCase.actual, testCase.expected) {
		fail(t, "TypesDao.FindWhere", testCase.expected, testCase.actual)
	}
}

func testTypesDao_FindWhere_MultipleResults(t *testing.T) {
	// Initialize test case's variables
	var testCase = typesDaoMultipleTestCase{
		query: "id < ? " +
			"ORDER BY id",
		params: nil,
		expected: []models.PokemonType{
			*newPokemonType(1, "normal", nil, "normal"),
			*newPokemonType(2, "fire", nil, "fire"),
			*newPokemonType(3, "fighting", nil, "fighting"),
		},
		actual: nil,
	}
	testCase.params = append(testCase.params, 4)

	// Get result
	testCase.actual = TYPES_DAO.FindWhere(testCase.query, testCase.params...)

	// Expected vs Actual
	if !reflect.DeepEqual(testCase.actual, testCase.expected) {
		fail(t, "TypesDao.FindWhere", testCase.expected, testCase.actual)
	}
}

func TestTypesDao_CRUD(t *testing.T) {
	testTypesDao_Create(t)
	testTypesDao_Save(t)
	testTypesDao_Delete(t)
	testTypesDao_FindOrCreate(t)
}

func testTypesDao_Create(t *testing.T) {
	// Set up test variables
	var (
		pokemonType models.PokemonType
		displayName = "Test"
		verify      = "SELECT id " +
			"FROM pvpgo.types " +
			"WHERE display_name = ?"
		expectedId int64
		cleanUp    = "DELETE FROM pvpgo.types " +
			"WHERE display_name = ?"
	)

	// Cleanup, even if error failure
	defer LIVE.Exec(cleanUp, displayName)

	// Create new type "Test"
	pokemonType = *TYPES_DAO.Create(displayName, nil, displayName)

	// Query the expected id
	CheckError(LIVE.QueryRow(verify, displayName).Scan(&expectedId))

	// Compare expected vs actual
	if pokemonType.Id() != expectedId {
		fail(t, "TypesDao.Create", expectedId, pokemonType.Id())
	}
}

func testTypesDao_Save(t *testing.T) {
	// Set up test variables
	var (
		pokemonType models.PokemonType
		displayName = "Test"
		insert      = "INSERT INTO pvpgo.types (first_type, displayName) " +
			"VALUES (?, ?)"
		id     int64
		verify = "SELECT display_name " +
			"FROM pvpgo.types " +
			"WHERE id = ?"
		expectedName = "Test/Test"
		actualName   string
		cleanUp      = "DELETE FROM pvpgo.types " +
			"WHERE id = ?"
		result sql.Result
		err    error
	)

	// Insert the test record into the DB
	result, err = LIVE.Exec(insert, displayName, displayName)
	CheckError(err)

	// Get our new record's ID
	id, err = result.LastInsertId()
	CheckError(err)

	// Cleanup, even if error failure
	defer LIVE.Exec(cleanUp, id)

	// Create our pokemon type
	pokemonType = *newPokemonType(id, displayName, nil, displayName)

	// Update the type
	pokemonType.SetSecondType(displayName)
	TYPES_DAO.Save(pokemonType)

	// Get our actual name after update
	CheckError(LIVE.QueryRow(verify, id).Scan(&actualName))

	// Check expected vs actual
	if expectedName != actualName {
		fail(t, "TypesDao.Save", expectedName, actualName)
	}
}

func testTypesDao_Delete(t *testing.T) {
	// TODO: test the delete function
}

func testTypesDao_FindOrCreate(t *testing.T) {
	// TODO: test the find or create function
}
