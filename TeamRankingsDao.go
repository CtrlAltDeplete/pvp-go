package main

import (
	"database/sql"
)

type TeamRankingsDao struct{}

func (dao *TeamRankingsDao) FindSingleWhere(query string, params ...interface{}) (error, *TeamRankingDto) {
	var (
		id                                int64
		cup                               string
		allyOneId, allyTwoId, allyThreeId int64
		score                             float64
		rows                              *sql.Rows
		err                               error
		count                             = 0
	)
	query = "SELECT * " +
		"FROM pvpgo.team_rankings " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		count++
		CheckError(rows.Scan(&id, &cup, &allyOneId, &allyTwoId, &allyThreeId, &score))
		if count > 1 {
			break
		}
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	if count == 0 {
		return NO_ROWS, nil
	} else if count == 1 {
		return nil, newTeamRanking(id, cup, allyOneId, allyTwoId, allyThreeId, score)
	} else {
		return MULTIPLE_ROWS, nil
	}
}

func (dao *TeamRankingsDao) FindWhere(query string, params ...interface{}) []TeamRankingDto {
	var (
		rankings                          []TeamRankingDto
		rows                              *sql.Rows
		err                               error
		id                                int64
		cup                               string
		allyOneId, allyTwoId, allyThreeId int64
		score                             float64
	)
	query = "SELECT * " +
		"FROM pvpgo.team_rankings " +
		"WHERE " + query
	rows, err = LIVE.Query(query, params...)
	CheckError(err)
	for rows.Next() {
		CheckError(rows.Scan(&id, &cup, &allyOneId, &allyTwoId, &allyThreeId, &score))
		rankings = append(rankings, *newTeamRanking(id, cup, allyOneId, allyTwoId, allyThreeId, score))
	}
	CheckError(rows.Err())
	CheckError(rows.Close())
	return rankings
}

func (dao *TeamRankingsDao) Create(cup string, allyOneId, allyTwoId, allyThreeId int64, score float64) (error, *TeamRankingDto) {
	var (
		result sql.Result
		err    error
		id     int64
		query  = "INSERT INTO pvpgo.team_rankings (cup, ally_one_id, ally_two_id, ally_three_id, score) " +
			"VALUES (?, ?, ?, ?, ?)"
	)
	result, err = LIVE.Exec(query, cup, allyOneId, allyTwoId, allyThreeId, score)
	if err != nil {
		return err, nil
	}
	id, err = result.LastInsertId()
	if err != nil {
		return err, nil
	}
	return nil, newTeamRanking(id, cup, allyOneId, allyTwoId, allyThreeId, score)
}

func newTeamRanking(id int64, cup string, allyOneId, allyTwoId, allyThreeId int64, score float64) *TeamRankingDto {
	var ranking = TeamRankingDto{}
	ranking.SetId(id)
	ranking.SetCup(cup)
	ranking.SetAllyOneId(allyOneId)
	ranking.SetAllyTwoId(allyTwoId)
	ranking.SetAllyThreeId(allyThreeId)
	ranking.SetScore(score)
	return &ranking
}
