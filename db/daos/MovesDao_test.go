package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"fmt"
	"reflect"
	"testing"
)

func TestMovesDao_FindSingleWhere(t *testing.T) {
	// Initialize test variables
	var (
		query    string
		typeId   int64
		moveId   int64
		result   sql.Result
		err      error
		expected *dtos.MoveDto
		actual   *dtos.MoveDto
		setupSql = "INSERT INTO pvpgo.moves (name, type_id, power, turns, energy) " +
			"VALUES (?, ?, 0, 0, 0)"
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakeMoveTeardown()

	// Test setup
	typeId = fakeTypeSetup()
	result, err = LIVE.Exec(setupSql, "Test1", typeId)
	CheckError(err)
	moveId, err = result.LastInsertId()
	CheckError(err)
	_, err = LIVE.Exec(setupSql, "Test2", typeId)
	CheckError(err)

	/* Happy Path Test */
	// Prepare test variables
	query = "id <= ? " +
		"AND id >= ?"
	expected = newMove(moveId, "Test1", typeId, 0, 0, 0, nil, nil,
		nil, nil)
	err, actual = MOVES_DAO.FindSingleWhere(query, moveId, moveId)
	CheckError(err)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "MoveDao.FindSingleWhere failed during the happy path", *expected, *actual)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	err, _ = MOVES_DAO.FindSingleWhere(query, 0)

	// Check expected vs actual
	if err != NO_ROWS {
		fail(t, "MoveDao.FindSingleWhere failed on no results", NO_ROWS, err)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id >= ?"
	err, _ = MOVES_DAO.FindSingleWhere(query, moveId)

	// Check expected vs actual
	if err != MULTIPLE_ROWS {
		fail(t, "MoveDao.FindSingleWhere failed on multiple results", MULTIPLE_ROWS, err)
	}
}

func TestMovesDao_FindWhere(t *testing.T) {
	// Initialize test variables
	var (
		query    string
		typeId   int64
		id       int64
		moveIds  []int64
		result   sql.Result
		err      error
		expected []dtos.MoveDto
		actual   []dtos.MoveDto
		setupSql = "INSERT INTO pvpgo.moves (name, type_id, power, turns, energy) " +
			"VALUES (?, ?, 0, 0, 0)"
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakeMoveTeardown()

	// Test setup
	typeId = fakeTypeSetup()
	for i := int64(1); i < 3; i++ {
		result, err = LIVE.Exec(setupSql, fmt.Sprintf("Test%d", i), typeId)
		CheckError(err)
		id, err = result.LastInsertId()
		CheckError(err)
		moveIds = append(moveIds, id)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	expected = []dtos.MoveDto{}
	actual = MOVES_DAO.FindWhere(query, 0)

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "MovesDao.FindWhere failed on no results", expected, actual)
	}

	/* Multiple Results Test */
	// Prepare test variables
	if moveIds == nil {
		panic("Expected non-nil moveIds")
	}
	query = "id IN (?, ?)"
	expected = []dtos.MoveDto{
		*newMove(moveIds[0], "Test1", typeId, 0, 0, 0, nil, nil,
			nil, nil),
		*newMove(moveIds[1], "Test2", typeId, 0, 0, 0, nil, nil,
			nil, nil),
	}
	actual = MOVES_DAO.FindWhere(query, moveIds[0], moveIds[1])

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "MovesDao.FindWhere failed on multiple results", expected, actual)
	}
}

func TestMovesDao_Create(t *testing.T) {
	// Initialize test variables
	var (
		typeId   int64
		err      error
		expected *dtos.MoveDto
		actual   *dtos.MoveDto
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakeMoveTeardown()

	// Prepare test variables
	typeId = fakeTypeSetup()
	err, actual = MOVES_DAO.Create("Test", typeId, 0, 0, 0, nil, nil,
		nil, nil)
	CheckError(err)
	expected = newMove(actual.Id(), "Test", typeId, 0, 0, 0, nil, nil,
		nil, nil)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "PokemonDao.Create failed", *expected, *actual)
	}
}

func TestMovesDao_Update(t *testing.T) {
	// Initialize test variables
	var (
		typeId      int64
		moveId      int64
		name        string
		power       int64
		turns       int64
		energy      int64
		probability sql.NullFloat64
		stage_delta sql.NullInt64
		stats       sql.NullString
		target      sql.NullString
		err         error
		expected    *dtos.MoveDto
		actual      *dtos.MoveDto
		verifySql   = "SELECT * " +
			"FROM pvpgo.moves " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakeMoveTeardown()

	// Prepare test variables
	typeId = fakeTypeSetup()
	moveId = fakeMoveSetup()
	expected = newMove(moveId, "Test", typeId, 0, 0, 0, nil, nil,
		nil, nil)
	MOVES_DAO.Update(*expected)
	err = LIVE.QueryRow(verifySql, moveId).Scan(&moveId, &name, &typeId, &power, &turns, &energy, &probability,
		&stage_delta, &stats, &target)
	CheckError(err)
	actual = newMove(moveId, name, typeId, power, turns, energy, probability, stage_delta,
		stats, target)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "MovesDao.Update failed", *expected, *actual)
	}
}

func TestMovesDao_Delete(t *testing.T) {
	// Initialize test variables
	var (
		moveId    int64
		move      *dtos.MoveDto
		expected  int64 = 0
		actual    int64
		verifySql = "SELECT COUNT(*) " +
			"FROM pvpgo.moves " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakeTypeTeardown()
	defer fakeMoveTeardown()

	// Prepare test variables
	moveId = fakeMoveSetup()
	move = newMove(moveId, "Test", fakeTypeSetup(), 0, 0, 0, nil, nil,
		nil, nil)
	MOVES_DAO.Delete(*move)
	CheckError(LIVE.QueryRow(verifySql, moveId).Scan(&actual))

	// Check expected vs actual
	if expected != actual {
		fail(t, "MovesDao.Delete failed", expected, actual)
	}
}
