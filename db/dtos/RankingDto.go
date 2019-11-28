package dtos

import (
	"database/sql"
	"log"
)

type RankingDto struct {
	id          int64
	cup         string
	pokemonId   int64
	moveSetId   int64
	pokemonRank sql.NullInt64
	moveSetRank float64
}

func (r *RankingDto) Id() int64 {
	return r.id
}

func (r *RankingDto) SetId(id int64) {
	r.id = id
}

func (r *RankingDto) Cup() string {
	return r.cup
}

func (r *RankingDto) SetCup(cup string) {
	r.cup = cup
}

func (r *RankingDto) PokemonId() int64 {
	return r.pokemonId
}

func (r *RankingDto) SetPokemonId(pokemonId int64) {
	r.pokemonId = pokemonId
}

func (r *RankingDto) MoveSetId() int64 {
	return r.moveSetId
}

func (r *RankingDto) SetMoveSetId(moveSetId int64) {
	r.moveSetId = moveSetId
}

func (r *RankingDto) PokemonRank() sql.NullInt64 {
	return r.pokemonRank
}

func (r *RankingDto) SetPokemonRank(pokemonRank interface{}) {
	switch t := pokemonRank.(type) {
	case int64:
		r.pokemonRank.Valid = true
		r.pokemonRank.Int64 = t
	case nil:
		r.pokemonRank.Valid = false
		r.pokemonRank.Int64 = 0
	case sql.NullInt64:
		if t.Valid {
			r.SetPokemonRank(t.Int64)
		} else {
			r.SetPokemonRank(nil)
		}
	default:
		log.Fatalf("Unkown type %T.", t)
	}
}

func (r *RankingDto) MoveSetRank() float64 {
	return r.moveSetRank
}

func (r *RankingDto) SetMoveSetRank(moveSetRank float64) {
	r.moveSetRank = moveSetRank
}
