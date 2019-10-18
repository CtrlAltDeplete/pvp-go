package models

type PokemonHasMove struct {
	id        int64
	pokemonId int64
	moveId    int64
	isLegacy  bool
}

func (phm *PokemonHasMove) Id() int64 {
	return phm.id
}

func (phm *PokemonHasMove) SetId(id int64) {
	phm.id = id
}

func (phm *PokemonHasMove) PokemonId() int64 {
	return phm.pokemonId
}

func (phm *PokemonHasMove) SetPokemonId(pokemonId int64) {
	phm.pokemonId = pokemonId
}

func (phm *PokemonHasMove) MoveId() int64 {
	return phm.moveId
}

func (phm *PokemonHasMove) SetMoveId(moveId int64) {
	phm.moveId = moveId
}

func (phm *PokemonHasMove) IsLegacy() bool {
	return phm.isLegacy
}

func (phm *PokemonHasMove) SetIsLegacy(isLegacy bool) {
	phm.isLegacy = isLegacy
}
