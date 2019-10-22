package models

type TypeMultiplier struct {
	id            int64
	receivingType int64
	actingType    int64
	multiplier    float64
}

func (tm *TypeMultiplier) Id() int64 {
	return tm.id
}

func (tm *TypeMultiplier) SetId(id int64) {
	tm.id = id
}

func (tm *TypeMultiplier) ReceivingType() int64 {
	return tm.receivingType
}

func (tm *TypeMultiplier) SetReceivingType(receivingType int64) {
	tm.receivingType = receivingType
}

func (tm *TypeMultiplier) ActingType() int64 {
	return tm.receivingType
}

func (tm *TypeMultiplier) SetActingType(actingType int64) {
	tm.actingType = actingType
}

func (tm *TypeMultiplier) Multiplier() float64 {
	return tm.multiplier
}

func (tm *TypeMultiplier) SetMultiplier(multiplier float64) {
	tm.multiplier = multiplier
}
