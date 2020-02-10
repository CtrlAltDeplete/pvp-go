package main

type CpMultiplierDto struct {
	id         int64
	level      float64
	multiplier float64
}

func (cpm *CpMultiplierDto) Id() int64 {
	return cpm.id
}

func (cpm *CpMultiplierDto) SetId(id int64) {
	cpm.id = id
}

func (cpm *CpMultiplierDto) Level() float64 {
	return cpm.level
}

func (cpm *CpMultiplierDto) SetLevel(level float64) {
	cpm.level = level
}

func (cpm *CpMultiplierDto) Multiplier() float64 {
	return cpm.multiplier
}

func (cpm *CpMultiplierDto) SetMultiplier(multiplier float64) {
	cpm.multiplier = multiplier
}
