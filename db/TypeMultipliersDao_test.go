package db

import (
	"PvP-Go/models"
	"database/sql"
	"reflect"
	"testing"
)

func TestTypeMultipliersDao_FindSingleWhere(t *testing.T) {
	// Initialize test variables
	var (
		query            string
		typeId           int64
		typeMultiplierId int64
		multiplier       = 1.0
		result           sql.Result
		err              error
		expected         *models.TypeMultiplier
		actual           *models.TypeMultiplier
		setupSql         = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
		teardownSql = "DELETE FROM pvpgo.type_multipliers " +
			"WHERE receiving_type = ? " +
			"OR acting_type = ?"
	)

	// Defer teardown
	defer func() {
		_, err = LIVE.Exec(teardownSql, typeId, typeId)
		fakeTypeTeardown()
		CheckError(err)
	}()

	// Test setup
	typeId = fakeTypeSetup()
	result, err = LIVE.Exec(setupSql, typeId, typeId, multiplier)
	CheckError(err)
	typeMultiplierId, err = result.LastInsertId()
	CheckError(err)
	_, err = LIVE.Exec(setupSql, typeId, 1, multiplier)
	CheckError(err)

	/* Happy Path Test */
	// Prepare test variables
	query = "id <= ? " +
		"AND id >= ?"
	expected = newTypeMultiplier(typeMultiplierId, typeId, typeId, multiplier)
	err, actual = TYPE_MULTIPLIER_DAO.FindSingleWhere(query, typeMultiplierId, typeMultiplierId)
	CheckError(err)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypeMultipliersDao.FindSingleWhere failed during the happy path", *expected, *actual)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	err, _ = TYPE_MULTIPLIER_DAO.FindSingleWhere(query, 0)

	// Check expected vs actual
	if err != NO_ROWS {
		fail(t, "TypeMultipliersDao.FindSingleWhere failed on no results", NO_ROWS, err)
	}

	/* Multiple Results Test */
	// Prepare test variables
	query = "id >= ?"
	err, _ = TYPE_MULTIPLIER_DAO.FindSingleWhere(query, typeMultiplierId)

	// Check expected vs actual
	if err != MULTIPLE_ROWS {
		fail(t, "TypeMultipliersDao.FindSingleWhere failed on multiple results", MULTIPLE_ROWS, err)
	}
}

func TestTypeMultipliersDao_FindWhere(t *testing.T) {
	// Initialize test variables
	var (
		query             string
		typeId            int64
		id                int64
		typeMultiplierIds []int64
		result            sql.Result
		err               error
		expected          []models.TypeMultiplier
		actual            []models.TypeMultiplier
		setupSql          = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
		teardownSql = "DELETE FROM pvpgo.type_multipliers " +
			"WHERE receiving_type = ? " +
			"OR acting_type = ?"
	)

	// Defer teardown
	defer func() {
		_, err = LIVE.Exec(teardownSql, typeId, typeId)
		fakeTypeTeardown()
		CheckError(err)
	}()

	// Test setup
	typeId = fakeTypeSetup()
	for i := int64(1); i < 3; i++ {
		result, err = LIVE.Exec(setupSql, typeId, i, 1.0)
		CheckError(err)
		id, err = result.LastInsertId()
		CheckError(err)
		typeMultiplierIds = append(typeMultiplierIds, id)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	expected = []models.TypeMultiplier{}
	actual = TYPE_MULTIPLIER_DAO.FindWhere(query, 0)

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "TypeMultipliersDao.FindWhere failed on no results", expected, actual)
	}

	/* Multiple Results Test */
	// Prepare test variables
	if typeMultiplierIds == nil {
		panic("Expected non-nil typeMultiplierIds")
	}
	query = "id IN (?, ?)"
	expected = []models.TypeMultiplier{
		*newTypeMultiplier(typeMultiplierIds[0], typeId, 1, 1.0),
		*newTypeMultiplier(typeMultiplierIds[1], typeId, 2, 1.0),
	}
	actual = TYPE_MULTIPLIER_DAO.FindWhere(query, typeMultiplierIds[0], typeMultiplierIds[1])

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "TypeMultipliersDao.FindWhere failed on multiple results", expected, actual)
	}
}

func TestTypeMultipliersDao_Create(t *testing.T) {
	// Initialize test variables
	var (
		typeId      int64
		multiplier  = 1.0
		err         error
		expected    *models.TypeMultiplier
		actual      *models.TypeMultiplier
		teardownSql = "DELETE FROM pvpgo.type_multipliers " +
			"WHERE receiving_type = ?"
	)

	// Defer teardown
	defer func() {
		_, err = LIVE.Exec(teardownSql, typeId)
		fakeTypeTeardown()
		CheckError(err)
	}()

	// Prepare test variables
	typeId = fakeTypeSetup()
	err, actual = TYPE_MULTIPLIER_DAO.Create(typeId, typeId, multiplier)
	CheckError(err)
	expected = newTypeMultiplier(actual.Id(), typeId, typeId, multiplier)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypeMultipliersDao.Create failed", *expected, *actual)
	}
}

func TestTypeMultipliersDao_Update(t *testing.T) {
	// Initialize test variables
	var (
		typeId           int64
		typeMultiplierId int64
		receivingType    int64
		actingType       int64
		multiplier       float64
		result           sql.Result
		err              error
		expected         *models.TypeMultiplier
		actual           *models.TypeMultiplier
		setupSql         = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
		verifySql = "SELECT * " +
			"FROM pvpgo.type_multipliers " +
			"WHERE id = ?"
		teardownSql = "DELETE FROM pvpgo.type_multipliers " +
			"WHERE receiving_type = ?"
	)

	// Defer teardown
	defer func() {
		_, err = LIVE.Exec(teardownSql, typeId)
		fakeTypeTeardown()
		CheckError(err)
	}()

	// Prepare test variables
	typeId = fakeTypeSetup()
	result, err = LIVE.Exec(setupSql, typeId, typeId, 1.0)
	CheckError(err)
	typeMultiplierId, err = result.LastInsertId()
	CheckError(err)
	expected = newTypeMultiplier(typeMultiplierId, typeId, typeId, 1.64)
	TYPE_MULTIPLIER_DAO.Update(*expected)
	err = LIVE.QueryRow(verifySql, typeMultiplierId).Scan(&typeMultiplierId, &receivingType, &actingType, &multiplier)
	CheckError(err)
	actual = newTypeMultiplier(typeMultiplierId, receivingType, actingType, multiplier)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "TypeMultipliersDao.Update failed", *expected, *actual)
	}
}

func TestTypeMultipliersDao_Delete(t *testing.T) {
	// TODO: Finish Delete test
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
