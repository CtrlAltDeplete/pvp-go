package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
)

type MatchUpsDao struct{}

func (dao *MatchUpsDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.MatchUpsDto) {
	var (
		id                                   int64
		cup                                  string
		matchUpType                          string
		pokemonId                            int64
		matchUpOne, matchUpTwo, matchUpThree int64
		rows                                 *sql.Rows
		err                                  error
		count                                = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.match_ups " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &cup, &matchUpType, &pokemonId, &matchUpOne, &matchUpTwo, &matchUpThree))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newMatchUp(id, cup, matchUpType, pokemonId, matchUpOne, matchUpTwo, matchUpThree)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *MatchUpsDao) FindWhere(query string, params ...interface{}) []dtos.MatchUpsDto {
	var (
		matchUps                             []dtos.MatchUpsDto
		rows                                 *sql.Rows
		err                                  error
		id                                   int64
		cup                                  string
		matchUpType                          string
		pokemonId                            int64
		matchUpOne, matchUpTwo, matchUpThree int64
	)
	query = "SELECT * " +
		"FROM pvpgo.match_ups " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &cup, &matchUpType, &pokemonId, &matchUpOne, &matchUpTwo, &matchUpThree))
		matchUps = append(matchUps, *newMatchUp(id, cup, matchUpType, pokemonId, matchUpOne, matchUpTwo, matchUpThree))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return matchUps
}

func (dao *MatchUpsDao) Create(cup, matchUpType string, pokemonId, matchUpOne, matchUpTwo, matchUpThree int64) (error, *dtos.MatchUpsDto) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.match_ups (cup, match_up_type, pokemon_id, match_up_one, match_up_two, match_up_three) " +
			"VALUES (?, ?, ?, ?, ?, ?)"
	)
	result, err = LIVE.Exec(query, cup, matchUpType, pokemonId, matchUpOne, matchUpTwo, matchUpThree)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newMatchUp(id, cup, matchUpType, pokemonId, matchUpOne, matchUpTwo, matchUpThree)
}

func newMatchUp(id int64, cup string, matchUpType string, pokemonId, matchUpOne, matchUpTwo, matchUpThree int64) *dtos.MatchUpsDto {
	var matchUp = dtos.MatchUpsDto{}
	matchUp.SetId(id)
	matchUp.SetCup(cup)
	matchUp.SetMatchUpType(matchUpType)
	matchUp.SetPokemonId(pokemonId)
	matchUp.SetMatchUpOneId(matchUpOne)
	matchUp.SetMatchUpTwoId(matchUpTwo)
	matchUp.SetMatchUpThreeId(matchUpThree)
	return &matchUp
}
