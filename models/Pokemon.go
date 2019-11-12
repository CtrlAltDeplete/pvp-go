package models

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"math"
)

type Pokemon struct {
	name           string
	baseStats      map[string]float64
	ivs            map[string]float64
	stats          map[string]float64
	statBuffs      map[string]float64
	typeIds        []int64
	cp             float64
	maxHp          float64
	hp             float64
	level          float64
	priority       int64
	fastMove       Move
	chargeMoves    []Move
	typeEfficacy   map[int64]float64
	dps            float64
	energy         float64
	coolDown       int64
	shields        int64
	hasActed       bool
	baitShields    bool
	farmEnergy     bool
	bestChargeMove *Move
}

func (pokemon *Pokemon) MaximizeStats() {
	var (
		cpms                                       = map[float64]float64{}
		level, atk, def, sta, cp, score, bestScore float64
		ivs, stats                                 map[string]float64
	)
	for _, cpDto := range daos.CP_DAO.FindAll() {
		cpms[cpDto.Level()] = cpDto.Multiplier()
	}
	pokemon.ivs = map[string]float64{}
	pokemon.stats = map[string]float64{}
	shouldContinue := func(ivs map[string]float64, level float64) bool {
		cp, _, score := CalculateCp(pokemon.baseStats, ivs, cpms[level])
		return cp < 1480 || score <= bestScore
	}
	shouldBreak := func(ivs map[string]float64, level float64) bool {
		cp, _, _ := CalculateCp(pokemon.baseStats, ivs, cpms[level])
		return cp > 1500
	}
	bestScore = 0
	ivs = map[string]float64{}
	for level = 1.0; level <= 40.0; level += 0.5 {
		ivs["atk"] = 0
		ivs["def"] = 0
		ivs["sta"] = 0
		if shouldBreak(ivs, level) {
			break
		}

		ivs["atk"] = 15
		ivs["def"] = 15
		ivs["sta"] = 15
		if level < 40.0 && shouldContinue(ivs, level) {
			continue
		}
		for atk = 0.0; atk <= 15.0; atk++ {
			ivs["atk"] = atk
			ivs["def"] = 0
			ivs["sta"] = 0
			if shouldBreak(ivs, level) {
				break
			}

			ivs["atk"] = atk
			ivs["def"] = 15
			ivs["sta"] = 15
			if atk < 15 && shouldContinue(ivs, level) {
				continue
			}
			for def = 0.0; def <= 15.0; def++ {
				ivs["atk"] = atk
				ivs["def"] = def
				ivs["sta"] = 0
				if shouldBreak(ivs, level) {
					break
				}

				ivs["atk"] = atk
				ivs["def"] = def
				ivs["sta"] = 15
				if def < 15 && shouldContinue(ivs, level) {
					continue
				}
				for sta = 0.0; sta <= 15.0; sta++ {
					ivs["atk"] = atk
					ivs["def"] = def
					ivs["sta"] = sta
					cp, stats, score = CalculateCp(pokemon.baseStats, ivs, cpms[level])

					if cp > 1500 {
						break
					}

					if bestScore < score {
						pokemon.ivs["atk"] = ivs["atk"]
						pokemon.ivs["def"] = ivs["def"]
						pokemon.ivs["sta"] = ivs["sta"]
						pokemon.level = level
						pokemon.stats["atk"] = stats["atk"]
						pokemon.stats["def"] = stats["def"]
						pokemon.stats["sta"] = stats["sta"]
						pokemon.cp = cp
						bestScore = score
					}
				}
			}
		}
	}
}

func (pokemon *Pokemon) GetStab(move *Move) float64 {
	for _, id := range pokemon.typeIds {
		if id == move.typeId {
			return 1.2
		}
	}
	return 1
}

func (pokemon *Pokemon) GetEfficacy(typeId int64) float64 {
	if efficacy, ok := pokemon.typeEfficacy[typeId]; ok {
		return efficacy
	}
	return 1.0
}

func (pokemon *Pokemon) Reset() {
	pokemon.hp = pokemon.maxHp
	pokemon.energy = 0
	pokemon.coolDown = 0
	pokemon.shields = 0
	pokemon.statBuffs = map[string]float64{
		"atk": 0,
		"def": 0,
	}
}

func (pokemon *Pokemon) ApplyStatBuffs(buffStages map[string]float64) {
	var maxBuffStages = 4.0
	pokemon.statBuffs["atk"] += buffStages["atk"]
	pokemon.statBuffs["atk"] = math.Min(maxBuffStages, math.Max(-maxBuffStages, pokemon.statBuffs["atk"]))
	pokemon.statBuffs["def"] += buffStages["def"]
	pokemon.statBuffs["def"] = math.Min(maxBuffStages, math.Max(-maxBuffStages, pokemon.statBuffs["def"]))
}

func (pokemon *Pokemon) GetAttack() float64 {
	return pokemon.GetStat(pokemon.stats["atk"], pokemon.statBuffs["atk"])
}

func (pokemon *Pokemon) GetDefense() float64 {
	return pokemon.GetStat(pokemon.stats["def"], pokemon.statBuffs["def"])
}

func (pokemon *Pokemon) GetStat(stat, buff float64) float64 {
	var (
		multiplier  float64
		buffDivisor = 4.0
	)
	if buff > 0 {
		multiplier = (buffDivisor + buff) / buffDivisor
	} else {
		multiplier = buffDivisor / (buffDivisor - buff)
	}

	return stat * multiplier
}

func (pokemon *Pokemon) SetBestMove(enemy Pokemon) {
	var bestDamage = 0.0
	var bestEnergy = 0.0
	for i := range pokemon.chargeMoves {
		chargeMove := pokemon.chargeMoves[i]
		damage := CalculateDamage(*pokemon, enemy, chargeMove)
		if damage > bestDamage || (damage == bestDamage && bestEnergy > -chargeMove.energy) {
			bestDamage = damage
			bestEnergy = -chargeMove.energy
			pokemon.bestChargeMove = &chargeMove
		}
	}
}

func NewPokemon(pokemonDto dtos.PokemonDto, fastMoveDto dtos.MoveDto, chargeMoveDtos []dtos.MoveDto) *Pokemon {
	var (
		p                  = Pokemon{}
		typesDto           *dtos.TypeDto
		firstTypeName      string
		secondTypeName     string
		typeMultiplierDtos []dtos.TypeMultiplierDto
		err                error
	)

	p.name = pokemonDto.Name()
	p.baseStats = map[string]float64{}
	p.baseStats["atk"] = pokemonDto.Atk()
	p.baseStats["def"] = pokemonDto.Def()
	p.baseStats["sta"] = pokemonDto.Sta()
	p.statBuffs = map[string]float64{}

	p.typeIds = []int64{}
	err, typesDto = daos.TYPES_DAO.FindSingleWhere("id = ?", pokemonDto.TypeId())
	daos.CheckError(err)
	if typesDto.SecondTypeNullable().Valid {
		firstTypeName = typesDto.FirstType()
		err, firstTypeDto := daos.TYPES_DAO.FindSingleByType(firstTypeName)
		daos.CheckError(err)
		p.typeIds = append(p.typeIds, firstTypeDto.Id())

		secondTypeName = typesDto.SecondType()
		err, secondTypeDto := daos.TYPES_DAO.FindSingleByType(secondTypeName)
		daos.CheckError(err)
		p.typeIds = append(p.typeIds, secondTypeDto.Id())
	} else {
		p.typeIds = append(p.typeIds, typesDto.Id())
	}

	p.cp = 0
	p.hp = 0
	p.level = 0
	p.priority = 0

	p.fastMove = *NewMove(fastMoveDto)
	p.chargeMoves = []Move{}
	for _, moveDto := range chargeMoveDtos {
		p.chargeMoves = append(p.chargeMoves, *NewMove(moveDto))
	}

	p.typeEfficacy = map[int64]float64{}
	typeMultiplierDtos = daos.TYPE_MULTIPLIER_DAO.FindAllByReceivingType(typesDto.Id())
	for _, dto := range typeMultiplierDtos {
		p.typeEfficacy[dto.ActingType()] = dto.Multiplier()
	}

	p.dps = 0
	p.energy = 0
	p.coolDown = 0
	p.shields = 0
	p.hasActed = false
	p.baitShields = false
	p.farmEnergy = false
	p.bestChargeMove = nil

	if pokemonDto.OptLevel() <= 1 {
		p.MaximizeStats()
		pokemonDto.SetOptAtk(p.ivs["atk"])
		pokemonDto.SetOptDef(p.ivs["def"])
		pokemonDto.SetOptSta(p.ivs["sta"])
		pokemonDto.SetOptLevel(p.level)
		daos.POKEMON_DAO.Update(pokemonDto)
	} else {
		p.ivs = map[string]float64{
			"atk": pokemonDto.OptAtk(),
			"def": pokemonDto.OptDef(),
			"sta": pokemonDto.OptSta(),
		}
		p.level = pokemonDto.OptLevel()
		err, cpmDto := daos.CP_DAO.FindByLevel(p.level)
		daos.CheckError(err)
		p.cp, p.stats, _ = CalculateCp(p.baseStats, p.ivs, cpmDto.Multiplier())
	}
	p.maxHp = p.stats["sta"]

	return &p
}

func CalculateCp(baseStats, ivs map[string]float64, cpm float64) (
	cp float64, stats map[string]float64, score float64) {
	var (
		atk, def, sta float64
	)
	atk = cpm * (baseStats["atk"] + ivs["atk"])
	def = cpm * (baseStats["def"] + ivs["def"])
	sta = cpm * (baseStats["sta"] + ivs["sta"])

	stats = map[string]float64{}
	stats["atk"] = atk
	stats["def"] = def
	stats["sta"] = math.Floor(sta)
	cp = math.Floor(atk * math.Pow(def, 0.5) * math.Pow(sta, 0.5) / 10)
	score = atk * def * math.Floor(sta)
	return
}
