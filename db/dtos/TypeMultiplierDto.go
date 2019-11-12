package dtos

type TypeMultiplierDto struct {
	id            int64
	receivingType int64
	actingType    int64
	multiplier    float64
}

func (tm *TypeMultiplierDto) Id() int64 {
	return tm.id
}

func (tm *TypeMultiplierDto) SetId(id int64) {
	tm.id = id
}

func (tm *TypeMultiplierDto) ReceivingType() int64 {
	return tm.receivingType
}

func (tm *TypeMultiplierDto) SetReceivingType(receivingType int64) {
	tm.receivingType = receivingType
}

func (tm *TypeMultiplierDto) ActingType() int64 {
	return tm.actingType
}

func (tm *TypeMultiplierDto) SetActingType(actingType int64) {
	tm.actingType = actingType
}

func (tm *TypeMultiplierDto) Multiplier() float64 {
	return tm.multiplier
}

func (tm *TypeMultiplierDto) SetMultiplier(multiplier float64) {
	tm.multiplier = multiplier
}
