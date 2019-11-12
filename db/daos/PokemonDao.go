package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"strings"
)

type PokemonDao struct{}

func (dao *PokemonDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.PokemonDto) {
	var (
		id                               int64
		gen                              int64
		name                             string
		typeId                           int64
		atk, def, sta                    float64
		dateAdd                          string
		legendary, pvpEligible           bool
		optLevel, optAtk, optDef, optSta float64
		rows                             *sql.Rows
		err                              error
		count                            = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.pokemon " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &gen, &name, &typeId, &atk, &def, &sta, &dateAdd, &legendary, &pvpEligible, &optLevel,
			&optAtk, &optDef, &optSta))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newPokemon(id, gen, name, typeId, atk, def, sta, dateAdd, legendary, pvpEligible, optLevel, optAtk,
			optDef, optSta)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *PokemonDao) FindById(id int64) (error, *dtos.PokemonDto) {
	var (
		query = "id = ?"
	)
	return dao.FindSingleWhere(query, id)
}

func (dao *PokemonDao) FindByName(name string) (error, *dtos.PokemonDto) {
	var (
		query = "name = ?"
	)
	return dao.FindSingleWhere(query, name)
}

func (dao *PokemonDao) FindWhere(query string, params ...interface{}) []dtos.PokemonDto {
	var (
		pokemon                          = []dtos.PokemonDto{}
		rows                             *sql.Rows
		e                                error
		id                               int64
		gen                              int64
		name                             string
		typeId                           int64
		atk, def, sta                    float64
		dateAdd                          string
		legendary, pvpEligible           bool
		optLevel, optAtk, optDef, optSta float64
	)
	query = "SELECT * " +
		"FROM pvpgo.pokemon " +
		"WHERE " + query
	rows, e = LIVE.Query(query, params...)
	CheckError(e)
	for rows.Next() {
		CheckError(rows.Scan(&id, &gen, &name, &typeId, &atk, &def, &sta, &dateAdd, &legendary, &pvpEligible,
			&optLevel, &optAtk, &optDef, &optSta))
		pokemon = append(pokemon, *newPokemon(id, gen, name, typeId, atk, def, sta, dateAdd, legendary, pvpEligible,
			optLevel, optAtk, optDef, optSta))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return pokemon
}

func (dao *PokemonDao) FindByGen(gen int64) []dtos.PokemonDto {
	var (
		query = "gen = ?"
	)
	return dao.FindWhere(query, gen)
}

func (dao *PokemonDao) FindByTypeId(id int64) []dtos.PokemonDto {
	var (
		query = "type_id = ?"
	)
	return dao.FindWhere(query, id)
}

func (dao *PokemonDao) FindByTypeIds(ids []int64) []dtos.PokemonDto {
	var (
		id     int64
		params []interface{}
		query  = "type_id INT (?" + strings.Repeat(", ?", len(ids)-1) + ")"
	)
	for _, id = range ids {
		params = append(params, id)
	}
	return dao.FindWhere(query, params...)
}

func (dao *PokemonDao) FindAll() []dtos.PokemonDto {
	return dao.FindWhere("TRUE")
}

func (dao *PokemonDao) Create(gen int64, name string, typeId int64, atk, def, sta float64, dateAdd string,
	legendary, pvpEligible bool, optLevel, optAtk, optDef, optSta float64) (error, *dtos.PokemonDto) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.pokemon (gen, name, type_id, atk, def, sta, date_add, is_legendary, is_pvp_eligible, opt_level, opt_atk, opt_def, opt_sta) " +
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	)
	result, err = LIVE.Exec(query, gen, name, typeId, atk, def, sta, dateAdd, legendary, pvpEligible,
		optLevel, optAtk, optDef, optSta)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newPokemon(id, gen, name, typeId, atk, def, sta, dateAdd, legendary, pvpEligible, optLevel, optAtk,
		optDef, optSta)
}

func (dao *PokemonDao) Update(pokemon dtos.PokemonDto) {
	var (
		err   error
		query = "UPDATE pvpgo.pokemon " +
			"SET gen = ?, " +
			"name = ?, " +
			"type_id = ?, " +
			"atk = ?, " +
			"def = ?, " +
			"sta = ?, " +
			"date_add = ?, " +
			"is_legendary = ?, " +
			"is_pvp_eligible = ?, " +
			"opt_level = ?, " +
			"opt_atk = ?, " +
			"opt_def = ?, " +
			"opt_sta = ? " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, pokemon.Gen(), pokemon.Name(), pokemon.TypeId(), pokemon.Atk(), pokemon.Def(),
		pokemon.Sta(), pokemon.DateAdd(), pokemon.Legendary(), pokemon.PvpEligible(), pokemon.OptLevel(),
		pokemon.OptAtk(), pokemon.OptDef(), pokemon.OptSta(), pokemon.Id())
	CheckError(err)
}

func (dao *PokemonDao) Upsert(gen int64, name string, typeId int64, atk, def, sta float64, dateAdd string,
	legendary, pvpEligible bool, optLevel, optAtk, optDef, optSta float64) (error, *dtos.PokemonDto) {
	var (
		err     error
		pokemon *dtos.PokemonDto
	)
	err, pokemon = dao.FindByName(name)
	if err == NO_ROWS {
		err, pokemon = dao.Create(gen, name, typeId, atk, def, sta, dateAdd, legendary, pvpEligible, optLevel,
			optAtk, optDef, optSta)
	} else if err == nil {
		pokemon.SetGen(gen)
		pokemon.SetTypeId(typeId)
		pokemon.SetAtk(atk)
		pokemon.SetDef(def)
		pokemon.SetSta(sta)
		pokemon.SetDateAdd(dateAdd)
		pokemon.SetLegendary(legendary)
		pokemon.SetPvpEligible(pvpEligible)
		pokemon.SetOptLevel(optLevel)
		pokemon.SetOptAtk(optAtk)
		pokemon.SetOptDef(optDef)
		pokemon.SetOptSta(sta)
		dao.Update(*pokemon)
	}
	if err != nil {
		return err, nil
	}
	return nil, pokemon
}

func (dao *PokemonDao) Delete(pokemon dtos.PokemonDto) {
	var (
		err   error
		query = "DELETE FROM pvpgo.pokemon " +
			"WHERE id = ?"
	)
	_, err = LIVE.Exec(query, pokemon.Id())
	CheckError(err)
}

func newPokemon(id int64, gen int64, name string, typeId int64, atk float64, def float64, sta float64, dateAdd string,
	legendary bool, pvpEligible bool, optLevel float64, optAtk float64, optDef float64, optSta float64) *dtos.PokemonDto {
	var p = dtos.PokemonDto{}
	p.SetId(id)
	p.SetGen(gen)
	p.SetName(name)
	p.SetTypeId(typeId)
	p.SetAtk(float64(atk))
	p.SetDef(float64(def))
	p.SetSta(float64(sta))
	p.SetDateAdd(dateAdd)
	p.SetLegendary(legendary)
	p.SetPvpEligible(pvpEligible)
	p.SetOptLevel(optLevel)
	p.SetOptAtk(float64(optAtk))
	p.SetOptDef(float64(optDef))
	p.SetOptSta(float64(optSta))
	return &p
}