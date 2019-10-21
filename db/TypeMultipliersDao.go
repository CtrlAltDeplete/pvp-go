package db

import "database/sql"

type TypeMultipliersDao struct{}

func (dao *TypeMultipliersDao) FindByIds(receivingType, actingType int64) (error, *float64) {
	var (
		multiplier float64
		query      = "SELECT multiplier " +
			"FROM pvpgo.type_multipliers " +
			"WHERE acting_type = ? " +
			"AND receiving_type = ?"
		rows  *sql.Rows
		err   error
		count = 0
	)
	rows, err = LIVE.Query(query, actingType, receivingType)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&multiplier))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else {
		return nil, &multiplier
	}
}

func (dao *TypeMultipliersDao) Create(receivingType, actingType int64, multiplier float64) float64 {
	var (
		e     error
		query = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
	)
	_, e = LIVE.Exec(query, receivingType, actingType, multiplier)
	CheckError(e)
	return multiplier
}

func (dao *TypeMultipliersDao) FindOrCreate(receivingType, actingType int64, multiplier float64) float64 {
	var (
		mult *float64
		err  error
	)
	err, mult = dao.FindByIds(receivingType, actingType)
	if err != nil {
		return dao.Create(receivingType, actingType, multiplier)
	}
	return *mult
}
