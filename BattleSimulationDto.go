package main

import "math"

type BattleSimulationDto struct {
	id                 int64
	allyId, enemyId    int64
	individualMatchups []int64
	score              int64
}

func (b *BattleSimulationDto) Id() int64 {
	return b.id
}

func (b *BattleSimulationDto) SetId(id int64) {
	b.id = id
}

func (b *BattleSimulationDto) AllyId() int64 {
	return b.allyId
}

func (b *BattleSimulationDto) SetAllyId(allyId int64) {
	b.allyId = allyId
}

func (b *BattleSimulationDto) EnemyId() int64 {
	return b.enemyId
}

func (b *BattleSimulationDto) SetEnemyId(enemyId int64) {
	b.enemyId = enemyId
}

func (b *BattleSimulationDto) IndividualMatchups() []int64 {
	return b.individualMatchups
}

func (b *BattleSimulationDto) SetIndividualMatchups(individualMatchups []int64) {
	b.individualMatchups = individualMatchups
	b.score = CalculateTotalScore(individualMatchups)
}

func (b *BattleSimulationDto) Score() float64 {
	return float64(b.score)
}

func CalculateTotalScore(individualMatchups []int64) int64 {
	return int64(math.Round(float64(individualMatchups[0]+2*individualMatchups[1]+individualMatchups[2]+
		2*individualMatchups[3]+4*individualMatchups[4]+2*individualMatchups[5]+
		individualMatchups[6]+2*individualMatchups[7]+individualMatchups[8]) / 16.0))
}
