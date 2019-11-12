package models

type Action struct {
	actionType string
	actor      int
	move       *Move
	turn       int64
	priority   int64
	valid      bool
	processed  bool
}

func NewAction(actionType string, actor int, move *Move, turn int64, priority int64) *Action {
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
