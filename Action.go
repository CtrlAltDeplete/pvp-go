package main

import "fmt"

var (
	CHARGE = "charge"
	FAST   = "fast"
)

type Action struct {
	actionType string
	actor      int
	move       Move
	turn       int64
	priority   int64
	valid      bool
	processed  bool
}

func (action *Action) ActionType() string {
	return action.actionType
}

func (action *Action) SetActionType(actionType string) {
	if actionType == CHARGE || actionType == FAST {
		action.actionType = actionType
	}
	panic(fmt.Sprintf("Invalid action type [%s]", actionType))
}

func (action *Action) IsCharge() bool {
	return action.actionType == CHARGE
}

func (action *Action) IsFast() bool {
	return action.actionType == FAST
}

func (action *Action) Actor() int {
	return action.actor
}

func (action *Action) SetActor(actorIndex int) {
	action.actor = actorIndex
}

func (action *Action) Enemy() int {
	return (action.actor + 1) % 2
}

func (action *Action) Move() Move {
	return action.move
}

func (action *Action) SetMove(move Move) {
	action.move = move
}

func (action *Action) Turn() int64 {
	return action.turn
}

func (action *Action) SetTurn(turn int64) {
	action.turn = turn
}

func (action *Action) Priority() int64 {
	return action.priority
}

func (action *Action) SetPriority(priority int64) {
	action.priority = priority
}

func (action *Action) Valid() bool {
	return action.valid
}

func (action *Action) SetValid(valid bool) {
	action.valid = valid
}

func (action *Action) Processed() bool {
	return action.processed
}

func (action *Action) SetProcessed(processed bool) {
	action.processed = processed
}

func NewAction(actionType string, actor int, move Move, turn int64, priority int64) *Action {
	var action = Action{}
	action.actionType = actionType
	action.actor = actor
	action.move = move
	action.turn = turn
	action.priority = priority
	action.valid = false
	action.processed = false
	return &action
}
