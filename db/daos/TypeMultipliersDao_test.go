package daos

import (
	"PvP-Go/db/dtos"
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
		expected         *dtos.TypeMultiplierDto
		actual           *dtos.TypeMultiplierDto
		setupSql         = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
	)

	// Defer teardown
	defer fakeTypeMultiplierTeardown()

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
		expected          []dtos.TypeMultiplierDto
		actual            []dtos.TypeMultiplierDto
		setupSql          = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
	)

	// Defer teardown
	defer fakeTypeMultiplierTeardown()

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
	expected = []dtos.TypeMultiplierDto{}
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
	expected = []dtos.TypeMultiplierDto{
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
		typeId     int64
		multiplier = 1.0
		err        error
		expected   *dtos.TypeMultiplierDto
		actual     *dtos.TypeMultiplierDto
	)

	// Defer teardown
	defer fakeTypeMultiplierTeardown()

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
		err              error
		expected         *dtos.TypeMultiplierDto
		actual           *dtos.TypeMultiplierDto
		verifySql        = "SELECT * " +
			"FROM pvpgo.type_multipliers " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakeTypeMultiplierTeardown()

	// Prepare test variables
	typeId, typeMultiplierId = fakeTypeMultiplierSetup()
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
	// Initialize test variables
	var (
		typeId           int64
		typeMultiplierId int64
		typeMultiplier   *dtos.TypeMultiplierDto
		expected         int64 = 0
		actual           int64
		verifySql        = "SELECT COUNT(*) " +
			"FROM pvpgo.type_multipliers " +
			"WHERE receiving_type = ?"
	)

	// Defer teardown
	defer fakeTypeMultiplierTeardown()

	// Prepare test variables
	typeId, typeMultiplierId = fakeTypeMultiplierSetup()
	typeMultiplier = newTypeMultiplier(typeMultiplierId, typeId, typeId, 1.0)
	TYPE_MULTIPLIER_DAO.Delete(*typeMultiplier)
	CheckError(LIVE.QueryRow(verifySql, typeId).Scan(&actual))

	// Check expected vs actual
	if expected != actual {
		fail(t, "TypeMultipliersDao.Delete failed", expected, actual)
	}
}
