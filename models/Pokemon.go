package models

type Pokemon struct {
	id                               int64
	gen                              int64
	name                             string
	typeId                           int64
	atk, def, sta                    float64
	dateAdd                          string
	legendary, pvpEligible           bool
	optLevel, optAtk, optDef, optSta float64
}

func (p *Pokemon) Id() int64 {
	return p.id
}

func (p *Pokemon) SetId(id int64) {
	p.id = id
}

func (p *Pokemon) Gen() int64 {
	return p.gen
}

func (p *Pokemon) SetGen(gen int64) {
	p.gen = gen
}

func (p *Pokemon) Name() string {
	return p.name
}

func (p *Pokemon) SetName(name string) {
	p.name = name
}

func (p *Pokemon) TypeId() int64 {
	return p.typeId
}

func (p *Pokemon) SetTypeId(typeId int64) {
	p.typeId = typeId
}

func (p *Pokemon) Atk() float64 {
	return p.atk
}

func (p *Pokemon) SetAtk(atk float64) {
	p.atk = atk
}

func (p *Pokemon) Def() float64 {
	return p.def
}

func (p *Pokemon) SetDef(def float64) {
	p.def = def
}

func (p *Pokemon) Sta() float64 {
	return p.sta
}

func (p *Pokemon) SetSta(sta float64) {
	p.sta = sta
}

func (p *Pokemon) DateAdd() string {
	return p.dateAdd
}

func (p *Pokemon) SetDateAdd(dateAdd string) {
	p.dateAdd = dateAdd
}

func (p *Pokemon) Legendary() bool {
	return p.legendary
}

func (p *Pokemon) SetLegendary(legendary bool) {
	p.legendary = legendary
}

func (p *Pokemon) PvpEligible() bool {
	return p.pvpEligible
}

func (p *Pokemon) SetPvpEligible(pvpEligible bool) {
	p.pvpEligible = pvpEligible
}

func (p *Pokemon) OptLevel() float64 {
	return p.optLevel
}

func (p *Pokemon) SetOptLevel(optLevel float64) {
	p.optLevel = optLevel
}

func (p *Pokemon) OptAtk() float64 {
	return p.optAtk
}

func (p *Pokemon) SetOptAtk(optAtk float64) {
	p.optAtk = optAtk
}

func (p *Pokemon) OptDef() float64 {
	return p.optDef
}

func (p *Pokemon) SetOptDef(optDef float64) {
	p.optDef = optDef
}

func (p *Pokemon) OptSta() float64 {
	return p.optSta
}

func (p *Pokemon) SetOptSta(optSta float64) {
	p.optSta = optSta
}
