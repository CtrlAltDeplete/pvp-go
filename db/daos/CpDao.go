package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
)

type CpDao struct{}

func (dao *CpDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.CpMultiplierDto) {
	var (
		id         int64
		level      float64
		multiplier float64
		rows       *sql.Rows
		err        error
		count      = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.cp_multipliers " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &level, &multiplier))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newCpMultiplier(id, level, multiplier)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *CpDao) FindByLevel(level float64) (error, *dtos.CpMultiplierDto) {
	var (
		query = "level = ?"
	)
	return dao.FindSingleWhere(query, level)
}

func (dao *CpDao) FindWhere(query string, params ...interface{}) []dtos.CpMultiplierDto {
	var (
		cpMultipliers = []dtos.CpMultiplierDto{}
		rows          *sql.Rows
		err           error
		id            int64
		level         float64
		multiplier    float64
	)
	query = "SELECT * " +
		"FROM pvpgo.cp_multipliers " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &level, &multiplier))
		cpMultipliers = append(cpMultipliers, *newCpMultiplier(id, level, multiplier))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return cpMultipliers
}

func (dao *CpDao) FindAll() []dtos.CpMultiplierDto {
	return dao.FindWhere("TRUE ORDER BY level ASC")
}

func (dao *CpDao) Create(level, multiplier float64) (error, *dtos.CpMultiplierDto) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.cp_multipliers (level, multiplier) " +
			"VALUES (?, ?)"
	)
	result, err = LIVE.Exec(query, level, multiplier)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newCpMultiplier(id, level, multiplier)
}

func (dao *CpDao) Update(cpMultiplier dtos.CpMultiplierDto) {
	var (
		err   error
		query = "UPDATE pvpgo.cp_multipliers " +
			"SET level = ?, " +
			"multiplier = ? " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, cpMultiplier.Level(), cpMultiplier.Multiplier(), cpMultiplier.Id())
	CheckError(err)
}

func (dao *CpDao) Upsert(level, multiplier float64) (error, *dtos.CpMultiplierDto) {
	var (
		err          error
		cpMultiplier *dtos.CpMultiplierDto
	)
	err, cpMultiplier = dao.FindByLevel(level)
	if err == NO_ROWS {
		err, cpMultiplier = dao.Create(level, multiplier)
	} else if err == nil {
		cpMultiplier.SetMultiplier(multiplier)
		dao.Update(*cpMultiplier)
	}
	if err != nil {
		return err, nil
	}
	return nil, cpMultiplier
}

func (dao *CpDao) Delete(cpMultiplier dtos.CpMultiplierDto) {
	var (
		err   error
		query = "DELETE FROM pvpgo.cp_multipliers " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, cpMultiplier.Id())
	CheckError(err)
}

func newCpMultiplier(id int64, level, multiplier float64) *dtos.CpMultiplierDto {
	var cpMult = dtos.CpMultiplierDto{}
	cpMult.SetId(id)
	cpMult.SetLevel(level)
	cpMult.SetMultiplier(multiplier)
	return &cpMult
}
