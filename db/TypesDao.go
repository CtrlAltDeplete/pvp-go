package db

import (
	"PvP-Go/models"
	"database/sql"
	"log"
	"strings"
)

type TypesDao struct{}

func (dao *TypesDao) FindSingleWhere(query string, params ...interface{}) (error, *models.PokemonType) {
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

func (dao *TypesDao) FindById(id int64) (error, *models.PokemonType) {
	var (
		query = "id = ?"
	)
	return dao.FindSingleWhere(query, id)
}

func (dao *TypesDao) FindSingleByType(t string) (error, *models.PokemonType) {
	var (
		query = "first_type = ? " +
			"AND second_type IS NULL"
	)
	return dao.FindSingleWhere(query, t)
}

func (dao *TypesDao) FindSingleByTypes(t1 string, t2 string) (error, *models.PokemonType) {
	var (
		query = "first_type IN (?, ?) " +
			"AND second_type IN (?, ?) " +
			"LIMIT 1"
	)
	return dao.FindSingleWhere(query, t1, t2, t1, t2)
}

func (dao *TypesDao) FindWhere(query string, params ...interface{}) []models.PokemonType {
	var (
		pokemonTypes = []models.PokemonType{}
		rows         *sql.Rows
		e            error
		id           int64
		firstType    string
		secondType   sql.NullString
		displayName  string
	)
	query = "SELECT * " +
		"FROM pvpgo.types " +
		"WHERE " + query
	rows, e = LIVE.Query(query, params...)
	CheckError(e)
	for rows.Next() {
		CheckError(rows.Scan(&id, &firstType, &secondType, &displayName))
		pokemonTypes = append(pokemonTypes, *newPokemonType(id, firstType, secondType, displayName))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return pokemonTypes
}

func (dao *TypesDao) FindAllByType(t string) []models.PokemonType {
	var (
		query = "first_type = ? " +
			"OR second_type = ?"
	)
	return dao.FindWhere(query, t, t)
}

func (dao *TypesDao) FindAllByTypes(t []string) []models.PokemonType {
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

func (dao *TypesDao) FindAll() []models.PokemonType {
	var (
		pokemonTypes = []models.PokemonType{}
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

func (dao *TypesDao) Create(firstType string, secondType interface{}, displayName string) *models.PokemonType {
	var (
		result sql.Result
		e      error
		id     int64
		query  = "INSERT INTO pvpgo.types (first_type, second_type, display_name) " +
			"VALUES (?, ?, ?)"
	)
	switch st := secondType.(type) {
	case string:
		result, e = LIVE.Exec(query, firstType, st, displayName)
	case nil:
		result, e = LIVE.Exec(query, firstType, sql.NullString{}, displayName)
	case sql.NullString:
		result, e = LIVE.Exec(query, firstType, st, displayName)
	default:
		log.Fatalf("Unknown type %T.", st)
	}
	CheckError(e)
	id, e = result.LastInsertId()
	CheckError(e)
	return newPokemonType(id, firstType, secondType, displayName)
}

func (dao *TypesDao) Save(pokemonType models.PokemonType) {
	var (
		e     error
		query = "UPDATE pvpgo.types " +
			"SET first_type = ?, " +
			"second_type = ?, " +
			"display_name = ? " +
			"WHERE id = ?"
	)
	_, e = LIVE.Exec(query, pokemonType.FirstType(), pokemonType.SecondType(), pokemonType.DisplayName(), pokemonType.Id())
	CheckError(e)
}

func (dao *TypesDao) Delete(pokemonType models.PokemonType) {
	var (
		e     error
		query = "DELETE FROM pvpgo.types " +
			"WHERE id = ?"
	)
	_, e = LIVE.Exec(query, pokemonType.Id())
	CheckError(e)
}

func (dao *TypesDao) FindOrCreate(types []string) *models.PokemonType {
	var (
		pokemonType *models.PokemonType
		err         error
	)
	if len(types) == 1 {
		err, pokemonType = dao.FindSingleByType(types[0])
		if err != nil {
			return dao.Create(types[0], nil, types[0])
		}
		return pokemonType
	} else if len(types) == 2 {
		err, pokemonType = dao.FindSingleByTypes(types[0], types[1])
		if err != nil {
			return dao.Create(types[0], types[1], types[0]+"/"+types[1])
		}
		return pokemonType
	}
	log.Fatalf("TypesDao.Upsert() expected 1 or 2 tyeps; got %d", len(types))
	return nil
}

func newPokemonType(id int64, firstType string, secondType interface{}, displayName string) *models.PokemonType {
	var pt = models.PokemonType{}
	pt.SetId(id)
	pt.SetFirstType(firstType)
	pt.SetSecondType(secondType)
	pt.SetDisplayName(displayName)
	return &pt
}
