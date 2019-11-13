package models

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"strings"
)

var (
	SELF     = "Self"
	OPPONENT = "Opponent"
)

type Move struct {
	id              int64
	name            string
	typeId          int64
	power           float64
	energy          float64
	coolDown        int64
	buffs           map[string]float64
	buffApplyChance float64
	buffTarget      string
}

func (move *Move) Id() int64 {
	return move.id
}

func (move *Move) SetId(id int64) {
	move.id = id
}

func (move *Move) Name() string {
	return move.name
}

func (move *Move) SetName(name string) {
	move.name = name
}

func (move *Move) TypeId() int64 {
	return move.typeId
}

func (move *Move) SetTypeId(typeId int64) {
	move.typeId = typeId
}

func (move *Move) Power() float64 {
	return move.power
}

func (move *Move) SetPower(power float64) {
	move.power = power
}

func (move *Move) Energy() float64 {
	return move.energy
}

func (move *Move) SetEnergy(energy float64) {
	move.energy = energy
}

func (move *Move) CoolDown() int64 {
	return move.coolDown
}

func (move *Move) SetCoolDown(coolDown int64) {
	move.coolDown = coolDown
}

func (move *Move) Buffs() map[string]float64 {
	return map[string]float64{
		ATK: move.buffs[ATK],
		DEF: move.buffs[DEF],
	}
}

func (move *Move) SetBuffs(buffs map[string]float64) {
	move.buffs[ATK] = buffs[ATK]
	move.buffs[DEF] = buffs[DEF]
}

func (move *Move) BuffApplyChance() float64 {
	return move.buffApplyChance
}

func (move *Move) SetBuffApplyChance(buffApplyChance float64) {
	move.buffApplyChance = buffApplyChance
}

func (move *Move) BuffTarget() string {
	return move.buffTarget
}

func (move *Move) SetBuffTarget(buffTarget string) {
	move.buffTarget = buffTarget
}

func (move *Move) DoesBuff() bool {
	return move.buffApplyChance == 1
}

func NewMove(moveDto dtos.MoveDto) *Move {
	var (
		buffTypes []string
		buffs     sql.NullString
	)
	move := Move{}
	move.id = moveDto.Id()
	move.name = moveDto.Name()
	move.typeId = moveDto.TypeId()
	move.power = moveDto.Power()
	move.energy = float64(moveDto.Energy())
	move.coolDown = moveDto.Turns()
	buffs = moveDto.StatsNullable()
	move.buffs = map[string]float64{
		ATK: 0,
		DEF: 0,
	}
	if buffs.Valid {
		buffTypes = strings.Split(buffs.String, ", ")
		for _, stat := range buffTypes {
			move.buffs[strings.ToLower(stat)] = float64(moveDto.StageDelta())
		}
		move.buffApplyChance = moveDto.Probability()
		move.buffTarget = moveDto.Target()
	} else {
		move.buffApplyChance = 0
		move.buffTarget = "Self"
	}
	return &move
}
