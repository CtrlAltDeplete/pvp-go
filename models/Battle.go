package models

import (
	"PvP-Go/db/dtos"
	"math"
	"math/rand"
	"reflect"
	"sort"
)

type Battle struct {
	pokemon                 []Pokemon
	turns                   int64
	lastProcessedTurn       int64
	queuedActions           []Action
	turnActions             []Action
	previousTurnActions     []Action
	roundChargeMovesUsed    int64
	roundChargeMovesStarted int64
	startingShields         int64
	chargeMoveUsed          bool
}

func (battle *Battle) Start() {
	battle.pokemon[0].Reset()
	battle.pokemon[1].Reset()

	battle.pokemon[0].SetBestMove(battle.pokemon[1])
	battle.pokemon[1].SetBestMove(battle.pokemon[0])

	battle.pokemon[0].shields = battle.startingShields
	battle.pokemon[1].shields = battle.startingShields

	battle.turns = 1
	battle.lastProcessedTurn = 0
	battle.queuedActions = []Action{}
	battle.turnActions = []Action{}
	battle.previousTurnActions = []Action{}
	battle.roundChargeMovesUsed = 0
	battle.roundChargeMovesStarted = 0
}

func (battle *Battle) Step() {
	battle.DecrementCooldowns()
	battle.RandomizePriority()
	battle.UpdateQueuedActions()

	sort.Slice(battle.turnActions, func(i, j int) bool {
		if battle.turnActions[i].turn == battle.turnActions[j].turn {
			return battle.turnActions[i].priority > battle.turnActions[j].priority
		}
		return battle.turnActions[i].turn < battle.turnActions[j].turn
	})

	battle.UpdateTurnActions()
}

func (battle *Battle) UpdateTurnActions() {
	for i := range battle.turnActions {
		action := battle.turnActions[i]
		poke := &battle.pokemon[action.actor]
		opponent := &battle.pokemon[(action.actor+1)%2]
		priorityChargeMoveThisTurn := false

		for _, otherAction := range battle.turnActions {
			if otherAction.actionType == "charge" {
				if otherAction.priority > 0 {
					priorityChargeMoveThisTurn = true
				}
			}
		}

		switch action.actionType {
		case "fast":
			action.valid = true
			if opponent.hp < 1 {
				action.valid = false
			}
			break
		case "charge":
			if poke.energy >= action.move.energy {
				action.valid = true
			}

			if poke.hp <= 0 && poke.priority == 0 && priorityChargeMoveThisTurn {
				action.valid = false
			}

			lethalFastMove := false
			opponentChargeMoveThisTurn := false

			for _, otherAction := range battle.turnActions {
				if action.actor != otherAction.actor {
					if otherAction.actionType == "fast" {
						if (opponent.coolDown == 0 && poke.hp <= CalculateDamage(
							battle.pokemon[otherAction.actor], battle.pokemon[(otherAction.actor+1)%2], *otherAction.move,
						)) || poke.hp < 1 {
							lethalFastMove = true
						}
					} else if otherAction.actionType == "charge" {
						opponentChargeMoveThisTurn = true
					}
				}
			}

			if lethalFastMove && !opponentChargeMoveThisTurn {
				action.valid = false
			}

			break
		}

		battle.ProcessAction(&action, poke, opponent)
	}

	battle.previousTurnActions = battle.turnActions
	battle.turnActions = []Action{}
	battle.lastProcessedTurn = battle.turns
	battle.turns++
}

func (battle *Battle) UpdateQueuedActions() {
	battle.roundChargeMovesUsed = 0
	actionsThisTurn := false
	chargeMoveThisTurn := false
	coolDownToSet := []int64{battle.pokemon[0].coolDown, battle.pokemon[1].coolDown}
	if battle.turns > battle.lastProcessedTurn {
		battle.turnActions = []Action{}
		for i := range battle.pokemon {
			battle.pokemon[i].hasActed = false
		}

		for i := range battle.pokemon {
			poke := &battle.pokemon[i]
			opponent := &battle.pokemon[(i+1)%2]
			action := battle.GetTurnAction(i, (i+1)%2)

			if action != nil {
				actionsThisTurn = true
				if action.actionType == "charge" {
					chargeMoveThisTurn = true
				} else if action.actionType == "fast" {
					coolDownToSet[i] += poke.fastMove.coolDown
				}

				if poke.hp > 0 && opponent.hp > 0 {
					battle.queuedActions = append(battle.queuedActions, *action)
				}
			}
		}
	}
	battle.pokemon[0].coolDown = coolDownToSet[0]
	battle.pokemon[1].coolDown = coolDownToSet[1]
	newQueuedActions := []Action{}
	for i := range battle.queuedActions {
		action := battle.queuedActions[i]
		valid := battle.IsActionValid(action, actionsThisTurn, chargeMoveThisTurn)

		if valid {
			battle.turnActions = append(battle.turnActions, action)
		} else {
			newQueuedActions = append(newQueuedActions, action)
		}
	}
	battle.queuedActions = newQueuedActions
}

func (battle *Battle) IsActionValid(action Action, actionsThisTurn bool, chargeMoveThisTurn bool) bool {
	valid := false
	if action.actionType == "fast" {
		turnsSinceActivated := battle.turns - action.turn
		chargeMoveLastTurn := false
		for _, previousAction := range battle.previousTurnActions {
			chargeMoveLastTurn = chargeMoveLastTurn || previousAction.actionType == "charge"
		}

		requiredTurnsToPass := battle.pokemon[action.actor].fastMove.coolDown - 1
		if turnsSinceActivated >= requiredTurnsToPass {
			action.priority += 20
			valid = true
		}
		if turnsSinceActivated >= 1 && chargeMoveLastTurn {
			action.priority += 20
			valid = true
		}

		if action.turn == battle.turns {
			chargeMoveLastTurn = false
			for _, previousAction := range battle.previousTurnActions {
				chargeMoveLastTurn = chargeMoveLastTurn || previousAction.actionType == "charge"
			}

			if actionsThisTurn {
				if turnsSinceActivated >= battle.pokemon[action.actor].fastMove.coolDown-1 {
					action.priority += 20
					valid = true
				}
				if turnsSinceActivated >= 1 && chargeMoveLastTurn {
					action.priority += 20
					valid = true
				}
			}

			chargeMoveLastTurn = false
			fastMoveRegisteredLastTurn := false

			for _, previousAction := range battle.previousTurnActions {
				if previousAction.actionType == "charge" && action.actor != previousAction.actor {
					chargeMoveLastTurn = true
				} else if previousAction.actionType == "fast" && action.actor == previousAction.actor {
					fastMoveRegisteredLastTurn = true
				}
			}

			if chargeMoveLastTurn && fastMoveRegisteredLastTurn && chargeMoveThisTurn {
				valid = false
			}
		}
	} else if action.actionType == "charge" {
		valid = true
	}
	return valid
}

func (battle *Battle) RandomizePriority() {
	if rand.Float64() > 0.5 {
		battle.pokemon[0].priority = 1
		battle.pokemon[1].priority = 0
	} else {
		battle.pokemon[0].priority = 0
		battle.pokemon[1].priority = 1
	}
}

func (battle *Battle) DecrementCooldowns() {
	for i := range battle.pokemon {
		battle.pokemon[i].coolDown--
		if battle.pokemon[i].coolDown < 0 {
			battle.pokemon[i].coolDown = 0
		}
	}
}

func (battle *Battle) GetTurnAction(pokeIndex, opponentIndex int) *Action {
	poke := &battle.pokemon[pokeIndex]
	opponent := &battle.pokemon[opponentIndex]
	var action *Action = nil

	battle.chargeMoveUsed = false

	if poke.coolDown == 0 && !poke.hasActed {
		poke.hasActed = true
		action = battle.DecideAction(pokeIndex, opponentIndex)

		if action == nil {
			action = NewAction("fast", pokeIndex, &poke.fastMove, battle.turns, poke.priority)
		} else {
			if action.actionType == "charged" {
				battle.roundChargeMovesStarted++
				if opponent.coolDown > 0 && !opponent.hasActed {
					action.priority += 4
					opponent.coolDown = 0

					opponentAction := battle.GetTurnAction(opponentIndex, pokeIndex)
					if opponentAction != nil && opponentAction.actionType == "charged" {
						battle.queuedActions = append(battle.queuedActions, *opponentAction)
					}
				}

				poke.coolDown = 0
				action.priority += 10
			}
		}
	}

	return action
}

func (battle *Battle) DecideAction(pokeIndex, opponentIndex int) *Action {
	var (
		poke, opponent                                           *Pokemon
		pokeFastDamage, pokeBestChargeDamage, opponentFastDamage float64
	)
	poke = &battle.pokemon[pokeIndex]
	opponent = &battle.pokemon[opponentIndex]
	pokeFastDamage = CalculateDamage(*poke, *opponent, poke.fastMove)
	pokeBestChargeDamage = CalculateDamage(*poke, *opponent, *poke.bestChargeMove)
	opponentFastDamage = CalculateDamage(*opponent, *poke, opponent.fastMove)

	if battle.ShouldUseBestCharge(poke, opponent, pokeFastDamage, pokeIndex) {
		return NewAction("charge", pokeIndex, poke.bestChargeMove, battle.turns, poke.priority)
	}

	for i := range poke.chargeMoves {
		move := &poke.chargeMoves[i]
		if battle.ShouldUseOtherChargeMove(poke, move, opponent, pokeIndex, pokeFastDamage, opponentFastDamage, pokeBestChargeDamage) {
			return NewAction("charge", pokeIndex, move, battle.turns, poke.priority)
		}
	}

	return nil
}

func (battle *Battle) ShouldUseOtherChargeMove(poke *Pokemon, move *Move, opponent *Pokemon, pokeIndex int, pokeFastDamage float64, opponentFastDamage float64, pokeBestChargeDamage float64) bool {
	if poke.energy >= -move.energy && !battle.chargeMoveUsed {
		damage := CalculateDamage(*poke, *opponent, *move)

		if damage >= opponent.hp && !battle.chargeMoveUsed {
			battle.chargeMoveUsed = true
			return true
		}

		if move.buffApplyChance == 1 && damage/(float64(poke.fastMove.coolDown)*opponent.stats["sta"]) >= 0.25 && opponent.hp > CalculateDamage(*poke, *opponent, *poke.bestChargeMove) {
			battle.chargeMoveUsed = true
			return true
		}

		if opponent.shields > 0 && !battle.chargeMoveUsed && (reflect.DeepEqual(move, poke.bestChargeMove) || poke.energy >= -poke.bestChargeMove.energy) {
			if opponent.hp > pokeFastDamage && opponent.hp > pokeFastDamage*float64(opponent.fastMove.coolDown)/float64(poke.fastMove.coolDown) {
				battle.chargeMoveUsed = true
				return true
			}
		}

		nearDeath := poke.hp <= opponentFastDamage && float64(opponent.coolDown)/float64(poke.fastMove.coolDown) < 3
		if poke.shields == 0 {
			for _, opponentChargeMove := range opponent.chargeMoves {
				if opponent.energy >= -opponentChargeMove.energy && poke.hp <= CalculateDamage(*opponent, *poke, opponentChargeMove) {
					nearDeath = true
				}
			}
		}

		if !nearDeath && ((((opponent.coolDown > 0) && (opponent.coolDown < poke.fastMove.coolDown)) || ((opponent.coolDown == 0) && (opponent.fastMove.coolDown < poke.fastMove.coolDown))) && (battle.roundChargeMovesUsed == 0)) {
			availableTurns := float64(poke.fastMove.coolDown - opponent.coolDown)
			futureActions := math.Ceil(availableTurns / float64(opponent.fastMove.coolDown))
			if opponent.fastMove.coolDown == 1 {
				futureActions++
			}
			if battle.roundChargeMovesUsed > 0 || battle.roundChargeMovesStarted > 0 {
				futureActions = 0
			}

			futureFastDamage := futureActions * opponentFastDamage

			if poke.hp <= futureFastDamage {
				nearDeath = true
			}

			if poke.shields == 0 {
				futureEffectiveEnergy := opponent.energy + opponent.fastMove.energy*(futureActions-1)
				futureEffectiveHp := poke.hp - ((futureActions - 1) * opponentFastDamage)

				if opponent.coolDown == 1 {
					futureEffectiveEnergy += opponent.fastMove.energy
				}

				for _, enemyMove := range opponent.chargeMoves {
					enemyMoveDamage := CalculateDamage(*opponent, *poke, enemyMove)
					if futureEffectiveEnergy >= -enemyMove.energy && futureEffectiveHp <= enemyMoveDamage {
						nearDeath = true
					}
				}
			}
		}

		if opponent.shields > 0 && opponent.hp <= pokeFastDamage {
			nearDeath = false
		}

		if poke.bestChargeMove != nil && poke.energy >= -poke.bestChargeMove.energy && damage < pokeBestChargeDamage {
			nearDeath = false
		}

		if nearDeath && !battle.chargeMoveUsed {
			battle.chargeMoveUsed = true
			return true
		}
	}
	return false
}

func (battle *Battle) ShouldUseBestCharge(poke *Pokemon, opponent *Pokemon, pokeFastDamage float64, pokeIndex int) bool {
	if poke.bestChargeMove != nil && poke.energy >= -poke.bestChargeMove.energy {
		useChargeMove := opponent.hp > pokeFastDamage && (opponent.shields == 0 || opponent.hp > pokeFastDamage*float64(opponent.fastMove.coolDown)/float64(poke.fastMove.coolDown))
		for _, chargeMove := range poke.chargeMoves {
			if poke.energy >= -chargeMove.energy && (-chargeMove.energy < -poke.bestChargeMove.energy || (poke.bestChargeMove.energy == -chargeMove.energy && chargeMove.buffApplyChance > 0)) {
				useChargeMove = useChargeMove && !(opponent.shields > 0 || opponent.hp <= CalculateDamage(*poke, *opponent, chargeMove))
			}
		}
		if useChargeMove {
			battle.chargeMoveUsed = true
			return true
		}
	}
	return false
}

func (battle *Battle) ProcessAction(action *Action, poke, opponent *Pokemon) {
	if action == nil || action.valid == false || action.processed {
		return
	}
	action.processed = true
	switch action.actionType {
	case "fast":
		move := poke.fastMove
		battle.UseMove(poke, opponent, &move)
		break
	case "charge":
		move := action.move
		if poke.energy >= -move.energy {
			battle.UseMove(poke, opponent, move)
			battle.chargeMoveUsed = true
			battle.roundChargeMovesUsed++
		}
		break
	}
}

func (battle *Battle) UseMove(poke, opponent *Pokemon, move *Move) {
	damage := CalculateDamage(*poke, *opponent, *move)
	if move.energy < 0 {
		poke.energy += move.energy

		if opponent.shields > 0 {
			useShield := true
			if move.buffApplyChance == 1 && (move.buffs["atk"] > 0 || move.buffs["def"] < 0) {
				useShield = false
				postMoveHp := opponent.hp - damage
				var currentBuffs = map[string]float64{}
				if move.buffTarget == dtos.BUFF_SELF {
					currentBuffs["atk"] = poke.statBuffs["atk"]
					currentBuffs["def"] = poke.statBuffs["def"]
					poke.ApplyStatBuffs(move.buffs)
				} else if move.buffTarget == dtos.DEBUFF_ENEMY {
					currentBuffs["atk"] = opponent.statBuffs["atk"]
					currentBuffs["def"] = opponent.statBuffs["def"]
					opponent.ApplyStatBuffs(move.buffs)
				}
				fastDamage := CalculateDamage(*poke, *opponent, poke.fastMove)
				fastAttackCount := math.Ceil(math.Max(-move.energy-poke.energy, 0)/poke.fastMove.energy) + 2
				fastAttackDamage := fastAttackCount * fastDamage
				cycleDamage := (fastAttackDamage + 1) * float64(opponent.shields)

				if postMoveHp <= cycleDamage {
					useShield = true
				}

				if move.buffTarget == dtos.BUFF_SELF {
					poke.statBuffs = currentBuffs
				} else if move.buffTarget == dtos.DEBUFF_ENEMY {
					opponent.statBuffs = currentBuffs
				}

				for _, chargeMove := range poke.chargeMoves {
					if poke.energy >= -chargeMove.energy {
						chargeDamage := CalculateDamage(*poke, *opponent, chargeMove)
						if chargeDamage >= opponent.hp {
							useShield = true
						}
					}
				}
			}

			if useShield {
				damage = 1
				opponent.shields--
			}
		}
	} else {
		poke.energy += move.energy
		if poke.energy > 100 {
			poke.energy = 100
		}
	}

	opponent.hp = math.Max(0, opponent.hp-damage)

	if move.buffApplyChance == 1 {
		if move.buffTarget == dtos.BUFF_SELF {
			poke.ApplyStatBuffs(move.buffs)
		} else if move.buffTarget == dtos.DEBUFF_ENEMY {
			opponent.ApplyStatBuffs(move.buffs)
		}
	}
}

func (battle *Battle) Simulate() (int64, int64) {
	battle.Start()
	continueBattle := true
	for continueBattle {
		battle.Step()
		continueBattle = battle.pokemon[0].hp > 0 && battle.pokemon[1].hp > 0
	}
	return battle.GetBattleRating()
}

func (battle *Battle) GetBattleRating() (int64, int64) {
	var (
		healthMultiplier = 12.0
		energyMultiplier = 2.0
		shieldMultiplier = 3.0
		healthScores     []float64
		energyScores     []float64
		shieldScores     []float64
		totalScores      []float64
		sum              float64
	)
	healthScores = append(healthScores, healthMultiplier-healthMultiplier*battle.pokemon[1].hp/battle.pokemon[1].maxHp,
		healthMultiplier-healthMultiplier*battle.pokemon[0].hp/battle.pokemon[0].maxHp)
	for i := range battle.pokemon {
		if battle.pokemon[i].hp > 0 {
			energyScores = append(energyScores, energyMultiplier*battle.pokemon[i].energy/-battle.pokemon[i].bestChargeMove.energy)
		} else {
			energyScores = append(energyScores, 0)
		}
	}
	if battle.startingShields > 0 {
		shieldScores = append(shieldScores, shieldMultiplier-shieldMultiplier*float64(battle.pokemon[1].shields)/float64(battle.startingShields),
			shieldMultiplier-shieldMultiplier*float64(battle.pokemon[0].shields)/float64(battle.startingShields))
	} else {
		shieldScores = append(shieldScores, 0, 0)
	}
	totalScores = []float64{healthScores[0] + energyScores[0] + shieldScores[0],
		healthScores[1] + energyScores[1] + shieldScores[1]}
	sum = totalScores[0] + totalScores[1]
	return int64(math.Round(1000.0 * totalScores[0] / sum)), int64(math.Round(1000.0 * totalScores[1] / sum))
}

func NewBattle(pokemon []Pokemon, startingShields int64) *Battle {
	var battle = Battle{}
	battle.pokemon = pokemon
	battle.turns = 0
	battle.lastProcessedTurn = 0
	battle.startingShields = startingShields
	for i := range battle.pokemon {
		battle.pokemon[i].shields = startingShields
	}
	battle.queuedActions = []Action{}
	battle.turnActions = []Action{}
	battle.previousTurnActions = []Action{}
	battle.roundChargeMovesUsed = 0
	battle.roundChargeMovesStarted = 0
	return &battle
}

func CalculateDamage(attacker, defender Pokemon, move Move) float64 {
	var (
		bonusMultiplier    = 1.3
		efficacyMultiplier = defender.GetEfficacy(move.typeId)
	)
	return math.Floor(move.power*attacker.GetStab(&move)*(attacker.GetAttack()/defender.GetDefense())*efficacyMultiplier*0.5*bonusMultiplier) + 1
}
