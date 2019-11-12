package models

import (
	"PvP-Go/db/dtos"
	"fmt"
	"reflect"
	"testing"
)

func TestPokemon_MaximizeStats(t *testing.T) {
	var (
		pokemonDto    dtos.PokemonDto
		moveDto       dtos.MoveDto
		pokemon       Pokemon
		expectedLevel float64 = 14
		expectedIvs           = map[string]float64{
			"atk": 0,
			"def": 13,
			"sta": 8,
		}
	)
	pokemonDto = dtos.PokemonDto{}
	pokemonDto.SetId(757)
	pokemonDto.SetGen(5)
	pokemonDto.SetName("Reshiram")
	pokemonDto.SetTypeId(71)
	pokemonDto.SetAtk(275)
	pokemonDto.SetDef(211)
	pokemonDto.SetSta(205)
	pokemonDto.SetDateAdd("2019-06-01")
	pokemonDto.SetLegendary(false)
	pokemonDto.SetPvpEligible(true)
	pokemonDto.SetOptLevel(0)
	pokemonDto.SetOptAtk(0)
	pokemonDto.SetOptDef(0)
	pokemonDto.SetOptSta(0)

	moveDto = dtos.MoveDto{}
	moveDto.SetId(8)
	moveDto.SetName("Dragon Breath")
	moveDto.SetTypeId(14)
	moveDto.SetPower(4)
	moveDto.SetTurns(1)
	moveDto.SetEnergy(3)
	moveDto.SetProbability(nil)
	moveDto.SetStageDelta(nil)
	moveDto.SetStats(nil)
	moveDto.SetTarget(nil)

	pokemon = *NewPokemon(pokemonDto, moveDto, []dtos.MoveDto{})

	if expectedLevel != pokemon.level || !reflect.DeepEqual(expectedIvs, pokemon.ivs) {
		fmt.Printf("Expected level %f; got level %f\n", expectedLevel, pokemon.level)
		fmt.Printf("Expected ivs %v; got ivs %v\n", expectedIvs, pokemon.ivs)
		t.Fail()
	}
}
