package dtos

type MoveSetDto struct {
	id                    int64
	pokemonId             int64
	fastMoveId            int64
	primaryChargeMoveId   int64
	secondaryChargeMoveId *int64
}

func (m *MoveSetDto) Id() int64 {
	return m.id
}

func (m *MoveSetDto) SetId(id int64) {
	m.id = id
}

func (m *MoveSetDto) PokemonId() int64 {
	return m.pokemonId
}

func (m *MoveSetDto) SetPokemonId(pokemonId int64) {
	m.pokemonId = pokemonId
}

func (m *MoveSetDto) FastMoveId() int64 {
	return m.fastMoveId
}

func (m *MoveSetDto) SetFastMoveId(fastMoveId int64) {
	m.fastMoveId = fastMoveId
}

func (m *MoveSetDto) PrimaryChargeMoveId() int64 {
	return m.primaryChargeMoveId
}

func (m *MoveSetDto) SetPrimaryChargeMoveId(primaryChargeMoveId int64) {
	m.primaryChargeMoveId = primaryChargeMoveId
}

func (m *MoveSetDto) SecondaryChargeMoveId() *int64 {
	return m.secondaryChargeMoveId
}

func (m *MoveSetDto) SetSecondaryChargeMoveId(secondaryChargeMoveId *int64) {
	m.secondaryChargeMoveId = secondaryChargeMoveId
}
