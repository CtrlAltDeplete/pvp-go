package dtos

import (
	"database/sql"
	"log"
	"strings"
)

type MoveDto struct {
	id            int64
	name          string
	typeId        int64
	power         float64
	turns, energy int64
	probability   sql.NullFloat64
	stageDelta    sql.NullInt64
	stats         sql.NullString
	target        sql.NullString
}

var (
	BUFF_SELF    = "Self"
	DEBUFF_ENEMY = "Opponent"
)

func (m *MoveDto) Id() int64 {
	return m.id
}

func (m *MoveDto) SetId(id int64) {
	m.id = id
}

func (m *MoveDto) Name() string {
	return m.name
}

func (m *MoveDto) SetName(name string) {
	m.name = name
}

func (m *MoveDto) TypeId() int64 {
	return m.typeId
}

func (m *MoveDto) SetTypeId(typeId int64) {
	m.typeId = typeId
}

func (m *MoveDto) Power() float64 {
	return m.power
}

func (m *MoveDto) SetPower(power float64) {
	m.power = power
}

func (m *MoveDto) Turns() int64 {
	return m.turns
}

func (m *MoveDto) SetTurns(turns int64) {
	m.turns = turns
}

func (m *MoveDto) Energy() int64 {
	return m.energy
}

func (m *MoveDto) SetEnergy(energy int64) {
	m.energy = energy
}

func (m *MoveDto) Probability() float64 {
	return m.probability.Float64
}

func (m *MoveDto) ProbabilityNullable() sql.NullFloat64 {
	return m.probability
}

func (m *MoveDto) SetProbability(probability interface{}) {
	switch p := probability.(type) {
	case float64:
		m.probability.Valid = true
		m.probability.Float64 = p
	case nil:
		m.probability.Valid = false
		m.probability.Float64 = 0
	case sql.NullFloat64:
		if p.Valid {
			m.SetProbability(p.Float64)
		} else {
			m.SetProbability(nil)
		}
	default:
		log.Fatalf("Unknown type %T.", p)
	}
}

func (m *MoveDto) StageDelta() int64 {
	return m.stageDelta.Int64
}

func (m *MoveDto) StageDeltaNullable() sql.NullInt64 {
	return m.stageDelta
}

func (m *MoveDto) SetStageDelta(stageDelta interface{}) {
	switch sd := stageDelta.(type) {
	case int64:
		m.stageDelta.Valid = true
		m.stageDelta.Int64 = sd
	case nil:
		m.stageDelta.Valid = false
		m.stageDelta.Int64 = 0
	case sql.NullInt64:
		if sd.Valid {
			m.SetStageDelta(sd.Int64)
		} else {
			m.SetStageDelta(nil)
		}
	default:
		log.Fatalf("Unknown type %T.", sd)
	}
}

func (m *MoveDto) Stats() []string {
	return strings.Split(m.stats.String, ", ")
}

func (m *MoveDto) StatsNullable() sql.NullString {
	return m.stats
}

func (m *MoveDto) SetStats(stats interface{}) {
	switch s := stats.(type) {
	case string:
		m.stats.Valid = true
		m.stats.String = s
	case []string:
		m.stats.Valid = true
		m.stats.String = strings.Join(s, ", ")
	case nil:
		m.stats.Valid = false
		m.stats.String = ""
	case sql.NullString:
		if s.Valid {
			m.SetStats(s.String)
		} else {
			m.SetStats(nil)
		}
	default:
		log.Fatalf("Unknown type %T.", s)
	}
}

func (m *MoveDto) Target() string {
	return m.target.String
}

func (m *MoveDto) TargetNullable() sql.NullString {
	return m.target
}

func (m *MoveDto) SetTarget(target interface{}) {
	switch t := target.(type) {
	case string:
		m.target.Valid = true
		m.target.String = t
	case nil:
		m.target.Valid = false
		m.target.String = ""
	case sql.NullString:
		if t.Valid {
			m.SetTarget(t.String)
		} else {
			m.SetTarget(nil)
		}
	default:
		log.Fatalf("Unknown type %T.", t)
	}
}
