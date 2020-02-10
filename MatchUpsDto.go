package main

type MatchUpsDto struct {
	id             int64
	cup            string
	matchUpType    string
	pokemonId      int64
	matchUpOneId   int64
	matchUpTwoId   int64
	matchUpThreeId int64
}

func (t *MatchUpsDto) Id() int64 {
	return t.id
}

func (t *MatchUpsDto) SetId(id int64) {
	t.id = id
}

func (t *MatchUpsDto) Cup() string {
	return t.cup
}

func (t *MatchUpsDto) SetCup(cup string) {
	t.cup = cup
}

func (t *MatchUpsDto) MatchUpType() string {
	return t.matchUpType
}

func (t *MatchUpsDto) SetMatchUpType(matchUpType string) {
	t.matchUpType = matchUpType
}

func (t *MatchUpsDto) PokemonId() int64 {
	return t.pokemonId
}

func (t *MatchUpsDto) SetPokemonId(pokemonId int64) {
	t.pokemonId = pokemonId
}

func (t *MatchUpsDto) MatchUpOneId() int64 {
	return t.matchUpOneId
}

func (t *MatchUpsDto) SetMatchUpOneId(matchUpOneId int64) {
	t.matchUpOneId = matchUpOneId
}

func (t *MatchUpsDto) MatchUpTwoId() int64 {
	return t.matchUpTwoId
}

func (t *MatchUpsDto) SetMatchUpTwoId(matchUpTwoId int64) {
	t.matchUpTwoId = matchUpTwoId
}

func (t *MatchUpsDto) MatchUpThreeId() int64 {
	return t.matchUpThreeId
}

func (t *MatchUpsDto) SetMatchUpThreeId(matchUpThreeId int64) {
	t.matchUpThreeId = matchUpThreeId
}
