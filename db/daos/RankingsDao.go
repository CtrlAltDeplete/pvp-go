package daos

import (
	"PvP-Go/db/dtos"
	"database/sql"
)

type RankingsDao struct{}

func (dao *RankingsDao) FindSingleWhere(query string, params ...interface{}) (error, *dtos.RankingDto) {
	var (
		id          int64
		cup         string
		pokemonId   int64
		moveSetId   int64
		pokemonRank sql.NullInt64
		moveSetRank float64
		rows        *sql.Rows
		err         error
		count       = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.rankings " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &cup, &pokemonId, &moveSetId, &pokemonRank, &moveSetRank))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newRanking(id, cup, pokemonId, moveSetId, pokemonRank, moveSetRank)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *RankingsDao) FindWhere(query string, params ...interface{}) []dtos.RankingDto {
	var (
		rankings    = []dtos.RankingDto{}
		rows        *sql.Rows
		err         error
		id          int64
		cup         string
		pokemonId   int64
		moveSetId   int64
		pokemonRank sql.NullInt64
		moveSetRank float64
	)
	query = "SELECT * " +
		"FROM pvpgo.rankings " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &cup, &pokemonId, &moveSetId, &pokemonRank, &moveSetRank))
		rankings = append(rankings, *newRanking(id, cup, pokemonId, moveSetId, pokemonRank, moveSetRank))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return rankings
}

func (dao *RankingsDao) Create(cup string, pokemonId int64, moveSetId int64, pokemonRank interface{},
	moveSetRank float64) (error, *dtos.RankingDto) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.rankings (cup, pokemon_id, move_set_id, pokemon_rank, move_set_rank) " +
			"VALUES (?, ?, ?, ?, ?)"
	)
	result, err = LIVE.Exec(query, cup, pokemonId, moveSetId, pokemonRank, moveSetRank)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newRanking(id, cup, pokemonId, moveSetId, pokemonRank, moveSetRank)
}

func newRanking(id int64, cup string, pokemonId int64, moveSetId int64, pokemonRank interface{}, moveSetRank float64) *dtos.RankingDto {
	var r = dtos.RankingDto{}
	r.SetId(id)
	r.SetCup(cup)
	r.SetPokemonId(pokemonId)
	r.SetMoveSetId(moveSetId)
	r.SetPokemonRank(pokemonRank)
	r.SetMoveSetRank(moveSetRank)
	return &r
}
