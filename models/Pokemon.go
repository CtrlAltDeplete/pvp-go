package models

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"math"
)

type Pokemon struct {
	pokemonDto       dtos.PokemonDto
	fastMove         dtos.MoveDto
	chargeMoves      []dtos.MoveDto
	level            float64
	ivs              map[string]float64
	stats            map[string]float64
	cpm              float64
	cp               int64
	maxHp, currentHp float64
	types            []string
}

func (pokemon *Pokemon) CalculateCp() int64 {
	var (
		atk, def, sta float64
		cpmDto        *dtos.CpMultiplierDto
		err           error
	)
	atk = pokemon.pokemonDto.Atk() + pokemon.ivs["atk"]
	def = math.Pow(pokemon.pokemonDto.Def()+pokemon.ivs["def"], 0.5)
	sta = math.Pow(pokemon.pokemonDto.Sta()+pokemon.ivs["sta"], 0.5)
	err, cpmDto = daos.CP_DAO.FindByLevel(pokemon.level)
	daos.CheckError(err)
	return CalculateCp(atk, def, sta, cpmDto.Multiplier())
}

func (pokemon *Pokemon) MaximizeStats() {
	type StatCombo struct {
		level, atk, def, sta, score, cpm float64
		cp                               int64
	}
	baseAtk := pokemon.pokemonDto.Atk()
	baseDef := pokemon.pokemonDto.Def()
	baseSta := pokemon.pokemonDto.Sta()
	bestOption := StatCombo{
		level: 0,
		atk:   0,
		def:   0,
		sta:   0,
		score: 0,
		cpm:   0,
		cp:    0,
	}
	for level := float64(1); level <= 40.0; level += 0.5 {
		for atkIv := float64(0); atkIv <= 15.0; atkIv++ {
			for defIv := float64(1); defIv <= 15.0; defIv++ {
				for staIv := float64(1); staIv <= 15.0; staIv++ {
					atk := baseAtk + atkIv
					def := math.Pow(baseDef+defIv, 0.5)
					sta := math.Pow(baseSta+staIv, 0.5)
					err, cpmDto := daos.CP_DAO.FindByLevel(level)
					daos.CheckError(err)
					cp := CalculateCp(atk, def, sta, cpmDto.Multiplier())
					if cp > 1500 {
						continue
					}
					score := atk * def * sta
					if score > bestOption.score {
						bestOption.atk = atkIv
						bestOption.def = defIv
						bestOption.sta = staIv
						bestOption.score = score
						bestOption.cpm = cpmDto.Multiplier()
						bestOption.cp = cp
					}
				}
			}
		}
	}
	pokemon.ivs["atk"] = bestOption.atk
	pokemon.ivs["def"] = bestOption.def
	pokemon.ivs["sta"] = bestOption.sta
	pokemon.cpm = bestOption.cpm
	pokemon.cp = bestOption.cp
	err, typeDto := daos.TYPES_DAO.FindById(pokemon.pokemonDto.TypeId())
	daos.CheckError(err)
	pokemon.types = []string{typeDto.FirstType()}
	if typeDto.SecondTypeNullable().Valid {
		pokemon.types = append(pokemon.types, typeDto.SecondType())
	}
}

func (pokemon *Pokemon) GetStab(move dtos.MoveDto) float64 {
	// TODO: Finsih
	return 0
}

func CalculateCp(atk, def, sta, cpm float64) int64 {
	return int64(math.Floor(atk*def*sta*math.Pow(cpm, 2)) / 10.0)
}

func InitializePokemon(pokemonDto dtos.PokemonDto, fastMove dtos.MoveDto, chargeMoves []dtos.MoveDto) *Pokemon {
	var pokemon = Pokemon{}
	pokemon.pokemonDto = pokemonDto
	pokemon.fastMove = fastMove
	pokemon.chargeMoves = chargeMoves
	pokemon.MaximizeStats()
	pokemon.stats["atk"] = pokemon.cpm * (pokemon.pokemonDto.Atk() + pokemon.ivs["atk"])
	pokemon.stats["def"] = pokemon.cpm * (pokemon.pokemonDto.Def() + pokemon.ivs["def"])
	pokemon.stats["sta"] = math.Max(math.Floor(pokemon.cpm*(pokemon.pokemonDto.Sta()+pokemon.ivs["sta"])), 10)
	pokemon.maxHp = pokemon.stats["sta"]
	pokemon.currentHp = pokemon.maxHp
	return &pokemon
}
