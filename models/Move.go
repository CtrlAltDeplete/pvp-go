package models

import (
	"PvP-Go/db/dtos"
	"database/sql"
	"strings"
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
		"atk": 0,
		"def": 0,
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
