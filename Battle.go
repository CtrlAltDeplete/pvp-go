package main

import (
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
	startingShields         []int64
	chargeMoveUsed          bool
}

func (battle *Battle) Start() {
	battle.pokemon[0].Reset()
	battle.pokemon[1].Reset()

	battle.pokemon[0].SetBestMove(battle.pokemon[1])
	battle.pokemon[1].SetBestMove(battle.pokemon[0])

	battle.pokemon[0].SetShields(battle.startingShields[0])
	battle.pokemon[1].SetShields(battle.startingShields[1])

	battle.turns = 1
	battle.lastProcessedTurn = 0
	battle.queuedActions = []Action{}
	battle.turnActions = []Action{}
	battle.previousTurnActions = []Action{}
	battle.roundChargeMovesUsed = 0
	battle.roundChargeMovesStarted = 0
}

func (battle *Battle) Step() {
	battle.DecrementCoolDowns()
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
		poke := &battle.pokemon[action.Actor()]
		opponent := &battle.pokemon[action.Enemy()]
		priorityChargeMoveThisTurn := false

		for _, otherAction := range battle.turnActions {
			if otherAction.IsCharge() {
				if otherAction.Priority() > 0 {
					priorityChargeMoveThisTurn = true
				}
			}
		}

		if action.IsFast() {
			action.valid = true
			if !opponent.IsAlive() {
				action.valid = false
			}
		} else if action.IsCharge() {
			move := action.Move()
			if poke.Energy() >= move.Energy() {
				action.SetValid(true)
			}

			if !opponent.IsAlive() && poke.Priority() == 0 && priorityChargeMoveThisTurn {
				action.SetValid(false)
			}

			lethalFastMove := false
			opponentChargeMoveThisTurn := false

			for _, otherAction := range battle.turnActions {
				if action.Actor() != otherAction.Actor() {
					if otherAction.IsFast() {
						if (opponent.CanAct() && poke.Hp() < CalculateDamage(
							battle.pokemon[otherAction.Actor()], battle.pokemon[otherAction.Enemy()], otherAction.Move(),
						)) || !poke.IsAlive() {
							lethalFastMove = true
						}
					} else if otherAction.IsCharge() {
						opponentChargeMoveThisTurn = true
					}
				}
			}

			if lethalFastMove && !opponentChargeMoveThisTurn {
				action.SetValid(false)
			}
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
	coolDownToSet := []int64{battle.pokemon[0].CoolDown(), battle.pokemon[1].CoolDown()}
	if battle.turns > battle.lastProcessedTurn {
		battle.turnActions = []Action{}
		for i := range battle.pokemon {
			battle.pokemon[i].SetHasActed(false)
		}

		for i := range battle.pokemon {
			poke := &battle.pokemon[i]
			opponent := &battle.pokemon[(i+1)%2]
			action := battle.GetTurnAction(i, (i+1)%2)

			if action != nil {
				actionsThisTurn = true
				if action.IsCharge() {
					chargeMoveThisTurn = true
				} else if action.IsFast() {
					coolDownToSet[i] += poke.fastMove.CoolDown()
				}

				if poke.IsAlive() && opponent.IsAlive() {
					battle.queuedActions = append(battle.queuedActions, *action)
				}
			}
		}
	}
	battle.pokemon[0].SetCoolDown(coolDownToSet[0])
	battle.pokemon[1].SetCoolDown(coolDownToSet[1])
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
	if action.IsFast() {
		turnsSinceActivated := battle.turns - action.Turn()
		chargeMoveLastTurn := false
		for _, previousAction := range battle.previousTurnActions {
			chargeMoveLastTurn = chargeMoveLastTurn || previousAction.IsCharge()
		}

		actorsFastMove := battle.pokemon[action.Actor()].FastMove()
		requiredTurnsToPass := actorsFastMove.CoolDown() - 1
		if turnsSinceActivated >= requiredTurnsToPass {
			action.SetPriority(action.Priority() + 20)
			valid = true
		}
		if turnsSinceActivated >= 1 && chargeMoveLastTurn {
			action.SetPriority(action.Priority() + 20)
			valid = true
		}

		if action.Turn() == battle.turns {
			chargeMoveLastTurn = false
			for _, previousAction := range battle.previousTurnActions {
				chargeMoveLastTurn = chargeMoveLastTurn || previousAction.IsCharge()
			}

			if actionsThisTurn {
				if turnsSinceActivated >= actorsFastMove.CoolDown()-1 {
					action.SetPriority(action.Priority() + 20)
					valid = true
				}
				if turnsSinceActivated >= 1 && chargeMoveLastTurn {
					action.SetPriority(action.Priority() + 20)
					valid = true
				}
			}

			chargeMoveLastTurn = false
			fastMoveRegisteredLastTurn := false

			for _, previousAction := range battle.previousTurnActions {
				if previousAction.IsCharge() && action.Actor() != previousAction.Actor() {
					chargeMoveLastTurn = true
				} else if previousAction.IsFast() && action.Actor() == previousAction.Actor() {
					fastMoveRegisteredLastTurn = true
				}
			}

			if chargeMoveLastTurn && fastMoveRegisteredLastTurn && chargeMoveThisTurn {
				valid = false
			}
		}
	} else if action.IsCharge() {
		valid = true
	}
	return valid
}

func (battle *Battle) RandomizePriority() {
	if rand.Float64() > 0.5 {
		battle.pokemon[0].SetPriority(1)
		battle.pokemon[1].SetPriority(0)
	} else {
		battle.pokemon[0].SetPriority(0)
		battle.pokemon[1].SetPriority(1)
	}
}

func (battle *Battle) DecrementCoolDowns() {
	for i := range battle.pokemon {
		battle.pokemon[i].DecrementCoolDown()
	}
}

func (battle *Battle) GetTurnAction(pokeIndex, opponentIndex int) *Action {
	poke := &battle.pokemon[pokeIndex]
	opponent := &battle.pokemon[opponentIndex]
	var action *Action = nil

	battle.chargeMoveUsed = false

	if poke.CanAct() && !poke.HasActed() {
		poke.SetHasActed(true)
		action = battle.DecideAction(pokeIndex, opponentIndex)

		if action == nil {
			action = NewAction(FAST, pokeIndex, poke.FastMove(), battle.turns, poke.Priority())
		} else {
			if action.IsCharge() {
				battle.roundChargeMovesStarted++
				if !opponent.CanAct() && !opponent.HasActed() {
					action.SetPriority(action.Priority() + 4)
					opponent.SetCoolDown(0)

					opponentAction := battle.GetTurnAction(opponentIndex, pokeIndex)
					if opponentAction != nil && opponentAction.IsCharge() {
						battle.queuedActions = append(battle.queuedActions, *opponentAction)
					}
				}

				poke.SetCoolDown(0)
				action.SetPriority(action.Priority() + 10)
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
	pokeFastDamage = CalculateDamage(*poke, *opponent, poke.FastMove())
	pokeBestChargeDamage = CalculateDamage(*poke, *opponent, poke.BestChargeMove())
	opponentFastDamage = CalculateDamage(*opponent, *poke, opponent.FastMove())

	if battle.ShouldUseBestCharge(poke, opponent, pokeFastDamage, pokeIndex) {
		return NewAction(CHARGE, pokeIndex, poke.BestChargeMove(), battle.turns, poke.Priority())
	}

	for i := range poke.chargeMoves {
		move := &poke.chargeMoves[i]
		if battle.ShouldUseOtherChargeMove(poke, move, opponent, pokeIndex, pokeFastDamage, opponentFastDamage, pokeBestChargeDamage) {
			return NewAction(CHARGE, pokeIndex, *move, battle.turns, poke.Priority())
		}
	}

	return nil
}

func (battle *Battle) ShouldUseOtherChargeMove(poke *Pokemon, move *Move, opponent *Pokemon, pokeIndex int, pokeFastDamage float64, opponentFastDamage float64, pokeBestChargeDamage float64) bool {
	pokeFastMove := poke.FastMove()
	opponentFastMove := opponent.FastMove()
	bestChargeMove := poke.BestChargeMove()
	if poke.Energy() >= -move.Energy() && !battle.chargeMoveUsed {
		damage := CalculateDamage(*poke, *opponent, *move)

		if damage >= opponent.Hp() && !battle.chargeMoveUsed {
			battle.chargeMoveUsed = true
			return true
		}

		if move.DoesBuff() && damage/(float64(pokeFastMove.CoolDown())*opponent.GetAttack()) >= 0.25 && opponent.Hp() > CalculateDamage(*poke, *opponent, bestChargeMove) {
			battle.chargeMoveUsed = true
			return true
		}

		if opponent.HasShields() && !battle.chargeMoveUsed && (reflect.DeepEqual(move, bestChargeMove) || poke.Energy() >= -bestChargeMove.Energy()) {
			if opponent.Hp() > pokeFastDamage && opponent.Hp() > pokeFastDamage*float64(pokeFastMove.CoolDown())/float64(pokeFastMove.CoolDown()) {
				battle.chargeMoveUsed = true
				return true
			}
		}

		nearDeath := poke.Hp() <= opponentFastDamage && float64(opponent.CoolDown())/float64(pokeFastMove.CoolDown()) < 3
		if !poke.HasShields() {
			for i := range opponent.ChargeMoves() {
				opponentChargeMove := opponent.ChargeMoves()[i]
				if opponent.Energy() >= -opponentChargeMove.Energy() && poke.Hp() <= CalculateDamage(*opponent, *poke, opponentChargeMove) {
					nearDeath = true
				}
			}
		}

		if !nearDeath && (((!opponent.CanAct() && opponent.CoolDown() < pokeFastMove.CoolDown()) || (opponent.CanAct() && opponentFastMove.CoolDown() < pokeFastMove.CoolDown())) && battle.roundChargeMovesUsed == 0) {
			availableTurns := float64(pokeFastMove.CoolDown() - opponent.CoolDown())
			futureActions := math.Ceil(availableTurns / float64(opponentFastMove.CoolDown()))
			if opponentFastMove.CoolDown() == 1 {
				futureActions++
			}
			if battle.roundChargeMovesUsed > 0 || battle.roundChargeMovesStarted > 0 {
				futureActions = 0
			}

			futureFastDamage := futureActions * opponentFastDamage

			if poke.Hp() <= futureFastDamage {
				nearDeath = true
			}

			if !poke.HasShields() {
				futureEffectiveEnergy := opponent.Energy() + opponentFastMove.Energy()*(futureActions-1)
				futureEffectiveHp := poke.Hp() - ((futureActions - 1) * opponentFastDamage)

				if opponent.CoolDown() == 1 {
					futureEffectiveEnergy += opponentFastMove.Energy()
				}

				for i := range opponent.ChargeMoves() {
					enemyMove := opponent.ChargeMoves()[i]
					enemyMoveDamage := CalculateDamage(*opponent, *poke, enemyMove)
					if futureEffectiveEnergy >= -enemyMove.Energy() && futureEffectiveHp <= enemyMoveDamage {
						nearDeath = true
					}
				}
			}
		}

		if opponent.HasShields() && opponent.Hp() <= pokeFastDamage {
			nearDeath = false
		}

		if poke.Energy() >= -bestChargeMove.Energy() && damage < pokeBestChargeDamage {
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
	if poke.energy >= -poke.bestChargeMove.Energy() {
		useChargeMove := opponent.hp > pokeFastDamage && (opponent.shields == 0 || opponent.hp > pokeFastDamage*float64(opponent.fastMove.CoolDown())/float64(poke.fastMove.CoolDown()))
		for _, chargeMove := range poke.chargeMoves {
			if poke.energy >= -chargeMove.Energy() && (-chargeMove.Energy() < -poke.bestChargeMove.Energy() || (poke.bestChargeMove.Energy() == -chargeMove.Energy() && chargeMove.DoesBuff())) {
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
	if action == nil || !action.Valid() || action.Processed() {
		return
	}
	action.SetProcessed(true)
	if action.IsFast() {
		move := poke.FastMove()
		battle.UseMove(poke, opponent, &move)
	} else if action.IsCharge() {
		move := action.Move()
		if poke.Energy() >= -move.Energy() {
			battle.UseMove(poke, opponent, &move)
			battle.chargeMoveUsed = true
			battle.roundChargeMovesUsed++
		}
	}
}

func (battle *Battle) UseMove(poke, opponent *Pokemon, move *Move) {
	damage := CalculateDamage(*poke, *opponent, *move)
	pokeFastMove := poke.FastMove()
	if move.Energy() < 0 {
		poke.SetEnergy(poke.Energy() + move.Energy())

		if opponent.HasShields() {
			useShield := true
			if move.DoesBuff() && (move.Buffs()[ATK] > 0 || move.Buffs()[DEF] < 0) {
				useShield = false
				postMoveHp := opponent.Hp() - damage
				var currentBuffs = map[string]float64{}
				if move.BuffTarget() == SELF {
					currentBuffs[ATK] = poke.StatBuffs()[ATK]
					currentBuffs[DEF] = poke.StatBuffs()[DEF]
					poke.ApplyStatBuffs(move.Buffs())
				} else if move.BuffTarget() == OPPONENT {
					currentBuffs[ATK] = opponent.StatBuffs()[ATK]
					currentBuffs[DEF] = opponent.StatBuffs()[DEF]
					opponent.ApplyStatBuffs(move.Buffs())
				}
				fastDamage := CalculateDamage(*poke, *opponent, pokeFastMove)
				fastAttackCount := math.Ceil(math.Max(-(move.Energy()+poke.Energy()), 0)/pokeFastMove.Energy()) + 2
				fastAttackDamage := fastAttackCount * fastDamage
				cycleDamage := (fastAttackDamage + 1) * float64(opponent.Shields())

				if postMoveHp <= cycleDamage {
					useShield = true
				}

				if move.BuffTarget() == SELF {
					poke.SetStatBuffs(currentBuffs)
				} else if move.BuffTarget() == OPPONENT {
					opponent.SetStatBuffs(currentBuffs)
				}

				for i := range poke.ChargeMoves() {
					chargeMove := poke.ChargeMoves()[i]
					if poke.Energy() >= -chargeMove.Energy() {
						chargeDamage := CalculateDamage(*poke, *opponent, chargeMove)
						if chargeDamage >= opponent.Hp() {
							useShield = true
						}
					}
				}
			}

			if useShield {
				damage = 1
				opponent.SetShields(opponent.Shields() - 1)
			}
		}
	} else {
		poke.SetEnergy(poke.Energy() + move.Energy())
	}

	opponent.SetHp(opponent.Hp() - damage)

	if move.DoesBuff() {
		if move.BuffTarget() == SELF {
			poke.ApplyStatBuffs(move.Buffs())
		} else if move.BuffTarget() == OPPONENT {
			opponent.ApplyStatBuffs(move.Buffs())
		}
	}
}

func (battle *Battle) Simulate() (int64, int64) {
	battle.Start()
	continueBattle := true
	for continueBattle {
		battle.Step()
		continueBattle = battle.pokemon[0].IsAlive() && battle.pokemon[1].IsAlive()
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

	healthScores = []float64{
		healthMultiplier - healthMultiplier*battle.pokemon[1].Hp()/battle.pokemon[1].MaxHp(),
		healthMultiplier - healthMultiplier*battle.pokemon[0].Hp()/battle.pokemon[0].MaxHp(),
	}

	energyScores = []float64{0, 0}
	for i := range battle.pokemon {
		if battle.pokemon[i].IsAlive() {
			bestChargeMove := battle.pokemon[i].BestChargeMove()
			energyScores[i] = energyMultiplier * battle.pokemon[i].Energy() / -bestChargeMove.Energy()
		}
	}

	shieldScores = []float64{0, 0}
	for i := range battle.pokemon {
		if battle.startingShields[i] > 0 {
			shieldScores[(i+1)%2] = shieldMultiplier - shieldMultiplier*float64(battle.pokemon[i].Shields())/float64(battle.startingShields[i])
		}
	}

	totalScores = []float64{
		healthScores[0] + energyScores[0] + shieldScores[0],
		healthScores[1] + energyScores[1] + shieldScores[1],
	}

	sum = totalScores[0] + totalScores[1]
	return int64(math.Round(1000.0 * totalScores[0] / sum)), int64(math.Round(1000.0 * totalScores[1] / sum))
}

func NewBattle(pokemon []Pokemon, startingShields []int64) *Battle {
	var battle = Battle{}
	battle.pokemon = pokemon
	battle.turns = 0
	battle.lastProcessedTurn = 0
	battle.startingShields = startingShields
	for i := range battle.pokemon {
		battle.pokemon[i].SetShields(startingShields[i])
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
		efficacyMultiplier = defender.GetEfficacy(move.TypeId())
	)
	return math.Floor(move.Power()*attacker.GetStab(&move)*(attacker.GetAttack()/defender.GetDefense())*efficacyMultiplier*0.5*bonusMultiplier) + 1
}

func DoAllBattles(pokemon []Pokemon) [][]int64 {
	var (
		allResults            = [][]int64{}
		allyShield            int64
		enemyShield           int64
		allyScore, enemyScore int64
		battle                Battle
	)
	allResults = [][]int64{}
	for allyShield = 0; allyShield < 3; allyShield++ {
		for enemyShield = 0; enemyShield < 3; enemyShield++ {
			battle = *NewBattle(pokemon, []int64{allyShield, enemyShield})
			allyScore, enemyScore = battle.Simulate()
			allResults = append(allResults, []int64{allyScore, enemyScore})
		}
	}
	return allResults
}
