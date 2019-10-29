package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"reflect"
	"testing"
)

func TestCpDao_FindSingleWhere(t *testing.T) {
	// Initialize test variables
	var (
		query          string
		cpMultiplierId int64
		result         sql.Result
		err            error
		expected       *dtos.CpMultiplierDto
		actual         *dtos.CpMultiplierDto
		setupSql       = "INSERT INTO pvpgo.cp_multipliers (level, multiplier) " +
			"VALUES (?, ?)"
	)

	// Defer teardown
	defer fakeCpMultiplierTeardown()

	// Test setup
	result, err = LIVE.Exec(setupSql, 0.0, 1.0)
	CheckError(err)
	cpMultiplierId, err = result.LastInsertId()
	CheckError(err)
	_, err = LIVE.Exec(setupSql, 0.5, 1.0)
	CheckError(err)

	/* Happy Path Test */
	// Prepare test variables
	query = "id <= ? " +
		"AND id >= ?"
	expected = newCpMultiplier(cpMultiplierId, 0.0, 1.0)
	err, actual = CP_DAO.FindSingleWhere(query, cpMultiplierId, cpMultiplierId)
	CheckError(err)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "CpDao.FindSingleWhere failed during the happy path", *expected, *actual)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	err, _ = CP_DAO.FindSingleWhere(query, 0)

	// Check expected vs actual
	if err != NO_ROWS {
		fail(t, "CpDao.FindSingleWhere failed on no results", NO_ROWS, err)
	}

	/* Multiple Results Test */
	// Prepare test variables
	query = "id >= ?"
	err, _ = CP_DAO.FindSingleWhere(query, cpMultiplierId)

	// Check expected vs actual
	if err != MULTIPLE_ROWS {
		fail(t, "CpDao.FindSingleWhere failed on multiple results", MULTIPLE_ROWS, err)
	}
}

func TestCpDao_FindWhere(t *testing.T) {
	// Initialize test variables
	var (
		query           string
		id              int64
		cpMultiplierIds []int64
		result          sql.Result
		err             error
		expected        []dtos.CpMultiplierDto
		actual          []dtos.CpMultiplierDto
		setupSql        = "INSERT INTO pvpgo.cp_multipliers (level, multiplier) " +
			"VALUES (?, ?)"
	)

	// Defer teardown
	defer fakeCpMultiplierTeardown()

	// Test setup
	for i := float64(0); i < 1; i += 0.5 {
		result, err = LIVE.Exec(setupSql, i, 1.0)
		CheckError(err)
		id, err = result.LastInsertId()
		CheckError(err)
		cpMultiplierIds = append(cpMultiplierIds, id)
	}

	/* No Results Test */
	// Prepare test variables
	query = "id <= ?"
	expected = []dtos.CpMultiplierDto{}
	actual = CP_DAO.FindWhere(query, 0)

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "CpDao.FindWhere failed on no results", expected, actual)
	}

	/* Multiple Results Test */
	// Prepare test variables
	if cpMultiplierIds == nil {
		panic("Expected non-nil cpMultiplierIds")
	}
	query = "id IN (?, ?)"
	expected = []dtos.CpMultiplierDto{
		*newCpMultiplier(cpMultiplierIds[0], float64(0.0), 1.0),
		*newCpMultiplier(cpMultiplierIds[1], float64(0.5), 1.0),
	}
	actual = CP_DAO.FindWhere(query, cpMultiplierIds[0], cpMultiplierIds[1])

	// Check expected vs actual
	if !reflect.DeepEqual(expected, actual) {
		fail(t, "CpDao.FindWhere failed on multiple results", expected, actual)
	}
}

func TestCpDao_Create(t *testing.T) {
	// Initialize test variables
	var (
		err      error
		expected *dtos.CpMultiplierDto
		actual   *dtos.CpMultiplierDto
	)

	// Defer teardown
	defer fakeCpMultiplierTeardown()

	// Prepare test variables
	err, actual = CP_DAO.Create(0.0, 1.0)
	CheckError(err)
	expected = newCpMultiplier(actual.Id(), 0.0, 1.0)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "CpDao.Create failed", *expected, *actual)
	}
}

func TestCpDao_Update(t *testing.T) {
	// Initialize test variables
	var (
		cpMultiplierId int64
		level          float64
		multiplier     float64
		err            error
		expected       *dtos.CpMultiplierDto
		actual         *dtos.CpMultiplierDto
		verifySql      = "SELECT * " +
			"FROM pvpgo.cp_multipliers " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakeCpMultiplierTeardown()

	// Prepare test variables
	cpMultiplierId = fakeCpMultiplierSetup()
	expected = newCpMultiplier(cpMultiplierId, 0.0, 1.0)
	CP_DAO.Update(*expected)
	err = LIVE.QueryRow(verifySql, cpMultiplierId).Scan(&cpMultiplierId, &level, &multiplier)
	CheckError(err)
	actual = newCpMultiplier(cpMultiplierId, level, multiplier)

	// Check expected vs actual
	if !reflect.DeepEqual(*expected, *actual) {
		fail(t, "CpDao.Update failed", *expected, *actual)
	}
}

func TestCpDao_Delete(t *testing.T) {
	// Initialize test variables
	var (
		cpMultiplierId int64
		cpMultiplier   *dtos.CpMultiplierDto
		expected       int64 = 0
		actual         int64
		verifySql      = "SELECT COUNT(*) " +
			"FROM pvpgo.cp_multipliers " +
			"WHERE id = ?"
	)

	// Defer teardown
	defer fakeCpMultiplierTeardown()

	// Prepare test variables
	cpMultiplierId = fakeCpMultiplierSetup()
	cpMultiplier = newCpMultiplier(cpMultiplierId, 0.0, 1.0)
	CP_DAO.Delete(*cpMultiplier)
	CheckError(LIVE.QueryRow(verifySql, cpMultiplierId).Scan(&actual))

	// Check expected vs actual
	if expected != actual {
		fail(t, "CpDao.Delete failed", expected, actual)
	}
}
