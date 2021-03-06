package main

import (
	"database/sql"
	"math"
	"sort"
	"strings"
)

type ApiDao struct{}

func (dao *ApiDao) GetRanking(cup string, pokemonId int64) *ApiIndividualRankingDto {
	var (
		query = `SELECT p.name, r.pokemon_rank, r.move_set_rank, r.move_set_id, fm.name, ft.display_name, pcm.name, pct.display_name, scm.name, sct.display_name
FROM pvpgo.rankings r
LEFT JOIN pvpgo.move_sets ms ON r.move_set_id = ms.id
LEFT JOIN pvpgo.pokemon p ON ms.pokemon_id = p.id
LEFT JOIN pvpgo.moves fm ON ms.fast_move_id = fm.id
LEFT JOIN pvpgo.types ft ON fm.type_id = ft.id
LEFT JOIN pvpgo.moves pcm ON ms.primary_charge_move_id = pcm.id
LEFT JOIN pvpgo.types pct ON pcm.type_id = pct.id
LEFT JOIN pvpgo.moves scm ON ms.secondary_charge_move_id = scm.id
LEFT JOIN pvpgo.types sct ON scm.type_id = sct.id
WHERE r.cup = ?
AND r.pokemon_id = ?
ORDER BY r.move_set_rank DESC`
		rows     *sql.Rows
		err      error
		response = ApiIndividualRankingDto{}
	)
	rows, err = LIVE.Query(query, cup, pokemonId)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return nil
	}
	for rows.Next() {
		var (
			name                    string
			pokemonRank             sql.NullInt64
			moveSetRank             float64
			moveSetId               int64
			fastMoveName            string
			fastMoveType            string
			primaryChargeMoveName   string
			primaryChargeMoveType   string
			secondaryChargeMoveName sql.NullString
			secondaryChargeMoveType sql.NullString
			fastMove                            = ApiMoveDto{}
			primaryChargeMove                   = ApiMoveDto{}
			secondaryChargeMove     *ApiMoveDto = nil
			moveSet                             = ApiMoveSetDto{}
		)
		err = rows.Scan(&name, &pokemonRank, &moveSetRank, &moveSetId, &fastMoveName, &fastMoveType,
			&primaryChargeMoveName, &primaryChargeMoveType, &secondaryChargeMoveName, &secondaryChargeMoveType)
		if err != nil {
			return nil
		}

		fastMove.Type = fastMoveType
		fastMove.Name = fastMoveName
		primaryChargeMove.Type = primaryChargeMoveType
		primaryChargeMove.Name = primaryChargeMoveName
		if secondaryChargeMoveType.Valid && secondaryChargeMoveName.Valid {
			secondaryChargeMove = &ApiMoveDto{Name: secondaryChargeMoveName.String, Type: secondaryChargeMoveType.String}
		}
		moveSet.Id = moveSetId
		moveSet.AbsoluteRank = moveSetRank
		moveSet.FastMove = fastMove
		moveSet.PrimaryChargeMove = primaryChargeMove
		moveSet.SecondaryChargeMove = secondaryChargeMove
		response.MoveSets = append(response.MoveSets, moveSet)

		if pokemonRank.Valid {
			response.Name = name
			response.RelativeRank = pokemonRank.Int64
		}
	}
	response.Id = pokemonId
	return &response
}

func (dao *ApiDao) GetAllRankingsForCup(cup string) *[]ApiCupRankingDto {
	var (
		query = `SELECT p.id, p.name, r.pokemon_rank, r.move_set_rank
FROM pvpgo.rankings r
LEFT JOIN pvpgo.pokemon p ON r.pokemon_id = p.id
WHERE r.cup = ?
AND r.pokemon_rank IS NOT NULL
ORDER BY r.move_set_rank DESC`
		rows        *sql.Rows
		err         error
		responseMap = map[int64]ApiCupRankingDto{}
	)
	rows, err = LIVE.Query(query, cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return nil
	}
	for rows.Next() {
		var (
			id          int64
			name        string
			pokemonRank sql.NullInt64
			moveSetRank float64
		)
		err = rows.Scan(&id, &name, &pokemonRank, &moveSetRank)
		if err != nil {
			return nil
		}

		if pokemonRank.Valid {
			response := ApiCupRankingDto{
				Name:         name,
				Id:           id,
				RelativeRank: pokemonRank.Int64,
				AbsoluteRank: math.Round(moveSetRank*10) / 10,
			}
			responseMap[id] = response
		} else {
			response := responseMap[id]
			responseMap[id] = response
		}
	}
	var response []ApiCupRankingDto
	for _, r := range responseMap {
		response = append(response, r)
	}
	sort.Slice(response, func(i, j int) bool {
		return response[i].RelativeRank < response[j].RelativeRank
	})
	return &response
}

func (dao *ApiDao) SaveCard(cup string, moveSetId int64) ([]byte, error) {
	err, moveSet := MOVE_SETS_DAO.FindSingleWhere("id = ?", moveSetId)
	if err != nil {
		return nil, err
	}
	rows, err := LIVE.Query("SELECT types FROM pvpgo.cups WHERE name = ?", cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return nil, err
	}
	var typesString string
	for rows.Next() {
		rows.Scan(&typesString)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	var types []string
	types = strings.Split(typesString, ",")
	return BuildCard(*moveSet, cup, types)
}
