package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"strings"
)

type TypesDao struct{}

func (dao *TypesDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.TypeDto) {
	var (
		id          int64
		firstType   string
		secondType  sql.NullString
		displayName string
		rows        *sql.Rows
		err         error
		count       = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.types " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &firstType, &secondType, &displayName))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newPokemonType(id, firstType, secondType, displayName)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *TypesDao) FindById(id int64) (error, *dtos.TypeDto) {
	var (
		query = "id = ?"
	)
	return dao.FindSingleWhere(query, id)
}

func (dao *TypesDao) FindSingleByType(t string) (error, *dtos.TypeDto) {
	var (
		query = "first_type = ? " +
			"AND second_type IS NULL"
	)
	return dao.FindSingleWhere(query, t)
}

func (dao *TypesDao) FindSingleByTypes(t []string) (error, *dtos.TypeDto) {
	if len(t) == 1 {
		return dao.FindSingleByType(t[0])
	} else if len(t) == 2 {
		var query = "first_type IN (?, ?) " +
			"AND second_type IN (?, ?) " +
			"LIMIT 1"
		return dao.FindSingleWhere(query, t[0], t[1], t[0], t[1])
	} else {
		return BAD_PARAMS, nil
	}
}

func (dao *TypesDao) FindWhere(query string, params ...interface{}) []dtos.TypeDto {
	var (
		pokemonTypes = []dtos.TypeDto{}
		rows         *sql.Rows
		err          error
		id           int64
		firstType    string
		secondType   sql.NullString
		displayName  string
	)
	query = "SELECT * " +
		"FROM pvpgo.types " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &firstType, &secondType, &displayName))
		pokemonTypes = append(pokemonTypes, *newPokemonType(id, firstType, secondType, displayName))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return pokemonTypes
}

func (dao *TypesDao) FindAllByType(t string) []dtos.TypeDto {
	var (
		query = "first_type = ? " +
			"OR second_type = ?"
	)
	return dao.FindWhere(query, t, t)
}

func (dao *TypesDao) FindAllByTypes(t []string) []dtos.TypeDto {
	var (
		params []interface{}
		query  = "first_type IN (?" + strings.Repeat(", ?", len(t)-1) + ") " +
			"OR second_type IN (?" + strings.Repeat(", ?", len(t)-1) + ")"
	)
	for i := 0; i < 2; i++ {
		for _, pt := range t {
			params = append(params, pt)
		}
	}
	return dao.FindWhere(query, params...)
}

func (dao *TypesDao) FindAll() []dtos.TypeDto {
	var (
		pokemonTypes = []dtos.TypeDto{}
		rows         *sql.Rows
		e            error
		id           int64
		firstType    string
		secondType   sql.NullString
		displayName  string
		query        = "SELECT * " +
			"FROM pvpgo.types"
	)
	rows, e = LIVE.Query(query)
	CheckError(e)
	for rows.Next() {
		CheckError(rows.Scan(&id, &firstType, &secondType, &displayName))
		pokemonTypes = append(pokemonTypes, *newPokemonType(id, firstType, secondType, displayName))
	}
	CheckError(rows.Close())
	return pokemonTypes
}

func (dao *TypesDao) Create(firstType string, secondType interface{}) (error, *dtos.TypeDto) {
	var (
		displayName string
		result      sql.Result
		err         error
		id          int64
		query       = "INSERT INTO pvpgo.types (first_type, second_type, display_name) " +
			"VALUES (?, ?, ?)"
	)
	switch st := secondType.(type) {
	case string:
		displayName = firstType + "/" + st
		result, err = LIVE.Exec(query, firstType, st, displayName)
	case nil:
		displayName = firstType
		result, err = LIVE.Exec(query, firstType, st, displayName)
	case sql.NullString:
		if st.Valid {
			return dao.Create(firstType, st.String)
		}
		return dao.Create(firstType, nil)
	default:
		result = nil
		err = BAD_PARAMS
	}
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newPokemonType(id, firstType, secondType, displayName)
}

func (dao *TypesDao) Update(pokemonType dtos.TypeDto) {
	var (
		err   error
		query = "UPDATE pvpgo.types " +
			"SET first_type = ?, " +
			"second_type = ?, " +
			"display_name = ? " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, pokemonType.FirstType(), pokemonType.SecondTypeNullable(), pokemonType.DisplayName(),
		pokemonType.Id())
	CheckError(err)
}

func (dao *TypesDao) Upsert(firstType string, secondType interface{}) (error, *dtos.TypeDto) {
	var (
		err         error
		pokemonType *dtos.TypeDto
		types       = []string{firstType}
	)
	switch st := secondType.(type) {
	case string:
		types = append(types, st)
	case sql.NullString:
		if st.Valid {
			return dao.Upsert(firstType, st.String)
		}
	}
	err, pokemonType = dao.FindSingleByTypes(types)
	if err == NO_ROWS {
		err, pokemonType = dao.Create(firstType, secondType)
	} else if err == nil {
		pokemonType.SetFirstType(firstType)
		pokemonType.SetSecondType(secondType)
		dao.Update(*pokemonType)
	}
	if err != nil {
		return err, nil
	}
	return nil, pokemonType
}

func (dao *TypesDao) Delete(pokemonType dtos.TypeDto) {
	var (
		err   error
		query = "DELETE FROM pvpgo.types " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, pokemonType.Id())
	CheckError(err)
}

func newPokemonType(id int64, firstType string, secondType interface{}, displayName string) *dtos.TypeDto {
	var pt = dtos.TypeDto{}
	pt.SetId(id)
	pt.SetFirstType(firstType)
	pt.SetSecondType(secondType)
	pt.SetDisplayName(displayName)
	return &pt
}
