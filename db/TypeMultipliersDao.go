package db

import (
	"PvP-Go/models"
	"database/sql"
)

type TypeMultipliersDao struct{}

func (dao *TypeMultipliersDao) FindSingleWhere(query string, params ...interface{}) (error, *models.TypeMultiplier) {
	var (
		id            int64
		receivingType int64
		actingType    int64
		multiplier    float64
		rows          *sql.Rows
		err           error
		count         = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.type_multipliers " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &receivingType, &actingType, &multiplier))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newTypeMultiplier(id, receivingType, actingType, multiplier)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *TypeMultipliersDao) FindByIds(receivingType, actingType int64) (error, *models.TypeMultiplier) {
	var (
		query = "receiving_type = ? " +
			"AND acting_type = ?"
	)
	return dao.FindSingleWhere(query, receivingType, actingType)
}

func (dao *TypeMultipliersDao) FindWhere(query string, params ...interface{}) []models.TypeMultiplier {
	var (
		typeMultipliers = []models.TypeMultiplier{}
		rows            *sql.Rows
		err             error
		id              int64
		receivingType   int64
		actingType      int64
		multiplier      float64
	)
	query = "SELECT * " +
		"FROM pvpgo.type_multipliers " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &receivingType, &actingType, &multiplier))
		typeMultipliers = append(typeMultipliers, *newTypeMultiplier(id, receivingType, actingType, multiplier))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return typeMultipliers
}

func (dao *TypeMultipliersDao) FindAllByReceivingType(receivingType int64) []models.TypeMultiplier {
	var (
		query = "receiving_type = ?"
	)
	return dao.FindWhere(query, receivingType)
}

func (dao *TypeMultipliersDao) FindAllByActingType(actingType int64) []models.TypeMultiplier {
	var (
		query = "acting_type = ?"
	)
	return dao.FindWhere(query, actingType)
}

func (dao *TypeMultipliersDao) Create(receivingType, actingType int64, multiplier float64) (error, *models.TypeMultiplier) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
	)
	result, err = LIVE.Exec(query, receivingType, actingType, multiplier)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newTypeMultiplier(id, receivingType, actingType, multiplier)
}

func (dao *TypeMultipliersDao) Update(typeMultiplier models.TypeMultiplier) {
	var (
		err   error
		query = "UPDATE pvpgo.type_multipliers " +
			"SET receiving_type = ?, " +
			"acting_type = ?, " +
			"multiplier = ? " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, typeMultiplier.ReceivingType(), typeMultiplier.ActingType(), typeMultiplier.Multiplier(),
		typeMultiplier.Id())
	CheckError(err)
}

func (dao *TypeMultipliersDao) Delete(typeMultiplier models.TypeMultiplier) {
	var (
		err   error
		query = "DELETE FROM pvpgo.type_multipliers " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, typeMultiplier.Id())
	CheckError(err)
}

func newTypeMultiplier(id int64, receivingType int64, actingType int64, multiplier float64) *models.TypeMultiplier {
	var tm = models.TypeMultiplier{}
	tm.SetId(id)
	tm.SetReceivingType(receivingType)
	tm.SetActingType(actingType)
	tm.SetMultiplier(multiplier)
	return &tm
}
