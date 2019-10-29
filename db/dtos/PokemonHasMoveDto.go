package dtos

type PokemonHasMoveDto struct {
	id        int64
	pokemonId int64
	moveId    int64
	isLegacy  bool
}

func (phm *PokemonHasMoveDto) Id() int64 {
	return phm.id
}

func (phm *PokemonHasMoveDto) SetId(id int64) {
	phm.id = id
}

func (phm *PokemonHasMoveDto) PokemonId() int64 {
	return phm.pokemonId
}

func (phm *PokemonHasMoveDto) SetPokemonId(pokemonId int64) {
	phm.pokemonId = pokemonId
}

func (phm *PokemonHasMoveDto) MoveId() int64 {
	return phm.moveId
}

func (phm *PokemonHasMoveDto) SetMoveId(moveId int64) {
	phm.moveId = moveId
}

func (phm *PokemonHasMoveDto) IsLegacy() bool {
	return phm.isLegacy
}

func (phm *PokemonHasMoveDto) SetIsLegacy(isLegacy bool) {
	phm.isLegacy = isLegacy
}
