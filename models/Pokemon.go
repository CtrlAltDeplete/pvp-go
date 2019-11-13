package models

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"math"
)

var (
	ATK = "atk"
	DEF = "def"
	STA = "sta"
)

type Pokemon struct {
	name           string
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
	energy         float64
	coolDown       int64
	shields        int64
	hasActed       bool
	bestChargeMove Move
}

func (pokemon *Pokemon) Name() string {
	return pokemon.name
}

func (pokemon *Pokemon) SetName(name string) {
	pokemon.name = name
}

func (pokemon *Pokemon) Stats() map[string]float64 {
	return map[string]float64{
		ATK: pokemon.stats[ATK],
		DEF: pokemon.stats[DEF],
		STA: pokemon.stats[STA],
	}
}

func (pokemon *Pokemon) SetStats(stats map[string]float64) {
	pokemon.stats[ATK] = stats[ATK]
	pokemon.stats[DEF] = stats[DEF]
	pokemon.stats[STA] = stats[STA]
	pokemon.maxHp = stats[STA]
}

func (pokemon *Pokemon) StatBuffs() map[string]float64 {
	return map[string]float64{
		ATK: pokemon.statBuffs[ATK],
		DEF: pokemon.statBuffs[DEF],
	}
}

func (pokemon *Pokemon) SetStatBuffs(statBuffs map[string]float64) {
	pokemon.statBuffs[ATK] = statBuffs[ATK]
	pokemon.statBuffs[DEF] = statBuffs[DEF]
}

func (pokemon *Pokemon) TypeIds() []int64 {
	return pokemon.typeIds
}

func (pokemon *Pokemon) SetTypeIds(typeIds []int64) {
	pokemon.typeIds = typeIds
}

func (pokemon *Pokemon) Cp() float64 {
	return pokemon.cp
}

func (pokemon *Pokemon) SetCp(cp float64) {
	pokemon.cp = cp
}

func (pokemon *Pokemon) MaxHp() float64 {
	return pokemon.maxHp
}

func (pokemon *Pokemon) SetMaxHp(maxHp float64) {
	pokemon.maxHp = maxHp
}

func (pokemon *Pokemon) Hp() float64 {
	return pokemon.hp
}

func (pokemon *Pokemon) IsAlive() bool {
	return pokemon.hp > 0
}

func (pokemon *Pokemon) SetHp(hp float64) {
	pokemon.hp = hp
	if pokemon.hp < 0 {
		pokemon.hp = 0
	}
}

func (pokemon *Pokemon) Level() float64 {
	return pokemon.level
}

func (pokemon *Pokemon) SetLevel(level float64) {
	pokemon.level = level
}

func (pokemon *Pokemon) Priority() int64 {
	return pokemon.priority
}

func (pokemon *Pokemon) SetPriority(priority int64) {
	pokemon.priority = priority
}

func (pokemon *Pokemon) FastMove() Move {
	return pokemon.fastMove
}

func (pokemon *Pokemon) SetFastMove(fastMove Move) {
	pokemon.fastMove = fastMove
}

func (pokemon *Pokemon) ChargeMoves() []Move {
	return pokemon.chargeMoves
}

func (pokemon *Pokemon) SetChargeMoves(chargeMoves []Move) {
	pokemon.chargeMoves = chargeMoves
}

func (pokemon *Pokemon) Energy() float64 {
	return pokemon.energy
}

func (pokemon *Pokemon) SetEnergy(energy float64) {
	pokemon.energy = energy
	if pokemon.energy > 100 {
		pokemon.energy = 100
	}
}

func (pokemon *Pokemon) CoolDown() int64 {
	return pokemon.coolDown
}

func (pokemon *Pokemon) CanAct() bool {
	return pokemon.coolDown == 0
}

func (pokemon *Pokemon) SetCoolDown(coolDown int64) {
	pokemon.coolDown = coolDown
}

func (pokemon *Pokemon) DecrementCoolDown() {
	pokemon.coolDown--
	if pokemon.coolDown < 0 {
		pokemon.coolDown = 0
	}
}

func (pokemon *Pokemon) Shields() int64 {
	return pokemon.shields
}

func (pokemon *Pokemon) HasShields() bool {
	return pokemon.shields > 0
}

func (pokemon *Pokemon) SetShields(shields int64) {
	pokemon.shields = shields
}

func (pokemon *Pokemon) HasActed() bool {
	return pokemon.hasActed
}

func (pokemon *Pokemon) SetHasActed(hasActed bool) {
	pokemon.hasActed = hasActed
}

func (pokemon *Pokemon) BestChargeMove() Move {
	return pokemon.bestChargeMove
}

func (pokemon *Pokemon) SetBestChargeMove(bestChargeMove Move) {
	pokemon.bestChargeMove = bestChargeMove
}

func (pokemon *Pokemon) MaximizeStats(baseStats map[string]float64) map[string]float64 {
	var (
		cpms                                       = map[float64]float64{}
		level, atk, def, sta, cp, score, bestScore float64
		ivs, bestIvs, stats                        map[string]float64
	)
	for _, cpDto := range daos.CP_DAO.FindAll() {
		cpms[cpDto.Level()] = cpDto.Multiplier()
	}
	pokemon.stats = map[string]float64{}
	shouldContinue := func(ivs map[string]float64, level float64) bool {
		cp, _, score := CalculateCp(baseStats, ivs, cpms[level])
		return cp < 1480 || score <= bestScore
	}
	shouldBreak := func(ivs map[string]float64, level float64) bool {
		cp, _, _ := CalculateCp(baseStats, ivs, cpms[level])
		return cp > 1500
	}
	bestScore = 0
	ivs = map[string]float64{}
	bestIvs = map[string]float64{}
	for level = 1.0; level <= 40.0; level += 0.5 {
		ivs[ATK] = 0
		ivs[DEF] = 0
		ivs[STA] = 0
		if shouldBreak(ivs, level) {
			break
		}

		ivs[ATK] = 15
		ivs[DEF] = 15
		ivs[STA] = 15
		if level < 40.0 && shouldContinue(ivs, level) {
			continue
		}
		for atk = 0.0; atk <= 15.0; atk++ {
			ivs[ATK] = atk
			ivs[DEF] = 0
			ivs[STA] = 0
			if shouldBreak(ivs, level) {
				break
			}

			ivs[ATK] = atk
			ivs[DEF] = 15
			ivs[STA] = 15
			if atk < 15 && shouldContinue(ivs, level) {
				continue
			}
			for def = 0.0; def <= 15.0; def++ {
				ivs[ATK] = atk
				ivs[DEF] = def
				ivs[STA] = 0
				if shouldBreak(ivs, level) {
					break
				}

				ivs[ATK] = atk
				ivs[DEF] = def
				ivs[STA] = 15
				if def < 15 && shouldContinue(ivs, level) {
					continue
				}
				for sta = 0.0; sta <= 15.0; sta++ {
					ivs[ATK] = atk
					ivs[DEF] = def
					ivs[STA] = sta
					cp, stats, score = CalculateCp(baseStats, ivs, cpms[level])

					if cp > 1500 {
						break
					}

					if bestScore < score {
						bestIvs[ATK] = ivs[ATK]
						bestIvs[DEF] = ivs[DEF]
						bestIvs[STA] = ivs[STA]
						pokemon.level = level
						pokemon.stats[ATK] = stats[ATK]
						pokemon.stats[DEF] = stats[DEF]
						pokemon.stats[STA] = stats[STA]
						pokemon.cp = cp
						bestScore = score
					}
				}
			}
		}
	}
	return bestIvs
}

func (pokemon *Pokemon) GetStab(move *Move) float64 {
	for _, id := range pokemon.typeIds {
		if id == move.TypeId() {
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
		ATK: 0,
		DEF: 0,
	}
}

func (pokemon *Pokemon) ApplyStatBuffs(buffStages map[string]float64) {
	var maxBuffStages = 4.0
	pokemon.statBuffs[ATK] += buffStages[ATK]
	pokemon.statBuffs[ATK] = math.Min(maxBuffStages, math.Max(-maxBuffStages, pokemon.statBuffs[ATK]))
	pokemon.statBuffs[DEF] += buffStages[DEF]
	pokemon.statBuffs[DEF] = math.Min(maxBuffStages, math.Max(-maxBuffStages, pokemon.statBuffs[DEF]))
}

func (pokemon *Pokemon) GetAttack() float64 {
	return pokemon.GetStat(pokemon.stats[ATK], pokemon.statBuffs[ATK])
}

func (pokemon *Pokemon) GetDefense() float64 {
	return pokemon.GetStat(pokemon.stats[DEF], pokemon.statBuffs[DEF])
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
		if damage > bestDamage || (damage == bestDamage && bestEnergy > -chargeMove.Energy()) {
			bestDamage = damage
			bestEnergy = -chargeMove.Energy()
			pokemon.bestChargeMove = chargeMove
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

	p.energy = 0
	p.coolDown = 0
	p.shields = 0
	p.hasActed = false

	baseStats := map[string]float64{
		ATK: pokemonDto.Atk(),
		DEF: pokemonDto.Def(),
		STA: pokemonDto.Sta(),
	}
	var ivs map[string]float64
	if pokemonDto.OptLevel() <= 1 {
		ivs = p.MaximizeStats(baseStats)
		pokemonDto.SetOptAtk(ivs[ATK])
		pokemonDto.SetOptDef(ivs[DEF])
		pokemonDto.SetOptSta(ivs[STA])
		pokemonDto.SetOptLevel(p.level)
		daos.POKEMON_DAO.Update(pokemonDto)
	} else {
		ivs := map[string]float64{
			ATK: pokemonDto.OptAtk(),
			DEF: pokemonDto.OptDef(),
			STA: pokemonDto.OptSta(),
		}
		p.level = pokemonDto.OptLevel()
		err, cpmDto := daos.CP_DAO.FindByLevel(p.level)
		daos.CheckError(err)
		p.cp, p.stats, _ = CalculateCp(baseStats, ivs, cpmDto.Multiplier())
	}
	p.maxHp = p.stats[STA]

	return &p
}

func CalculateCp(baseStats, ivs map[string]float64, cpm float64) (
	cp float64, stats map[string]float64, score float64) {
	var (
		atk, def, sta float64
	)
	atk = cpm * (baseStats[ATK] + ivs[ATK])
	def = cpm * (baseStats[DEF] + ivs[DEF])
	sta = cpm * (baseStats[STA] + ivs[STA])

	stats = map[string]float64{}
	stats[ATK] = atk
	stats[DEF] = def
	stats[STA] = math.Floor(sta)
	cp = math.Floor(atk * math.Pow(def, 0.5) * math.Pow(sta, 0.5) / 10)
	score = atk * def * math.Floor(sta)
	return
}
