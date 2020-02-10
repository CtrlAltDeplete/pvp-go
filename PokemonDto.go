package main

type PokemonDto struct {
	id                               int64
	gen                              int64
	name                             string
	typeId                           int64
	atk, def, sta                    float64
	dateAdd                          string
	legendary, pvpEligible           bool
	optLevel, optAtk, optDef, optSta float64
}

func (p *PokemonDto) Id() int64 {
	return p.id
}

func (p *PokemonDto) SetId(id int64) {
	p.id = id
}

func (p *PokemonDto) Gen() int64 {
	return p.gen
}

func (p *PokemonDto) SetGen(gen int64) {
	p.gen = gen
}

func (p *PokemonDto) Name() string {
	return p.name
}

func (p *PokemonDto) SetName(name string) {
	p.name = name
}

func (p *PokemonDto) TypeId() int64 {
	return p.typeId
}

func (p *PokemonDto) SetTypeId(typeId int64) {
	p.typeId = typeId
}

func (p *PokemonDto) Atk() float64 {
	return p.atk
}

func (p *PokemonDto) SetAtk(atk float64) {
	p.atk = atk
}

func (p *PokemonDto) Def() float64 {
	return p.def
}

func (p *PokemonDto) SetDef(def float64) {
	p.def = def
}

func (p *PokemonDto) Sta() float64 {
	return p.sta
}

func (p *PokemonDto) SetSta(sta float64) {
	p.sta = sta
}

func (p *PokemonDto) DateAdd() string {
	return p.dateAdd
}

func (p *PokemonDto) SetDateAdd(dateAdd string) {
	p.dateAdd = dateAdd
}

func (p *PokemonDto) Legendary() bool {
	return p.legendary
}

func (p *PokemonDto) SetLegendary(legendary bool) {
	p.legendary = legendary
}

func (p *PokemonDto) PvpEligible() bool {
	return p.pvpEligible
}

func (p *PokemonDto) SetPvpEligible(pvpEligible bool) {
	p.pvpEligible = pvpEligible
}

func (p *PokemonDto) OptLevel() float64 {
	return p.optLevel
}

func (p *PokemonDto) SetOptLevel(optLevel float64) {
	p.optLevel = optLevel
}

func (p *PokemonDto) OptAtk() float64 {
	return p.optAtk
}

func (p *PokemonDto) SetOptAtk(optAtk float64) {
	p.optAtk = optAtk
}

func (p *PokemonDto) OptDef() float64 {
	return p.optDef
}

func (p *PokemonDto) SetOptDef(optDef float64) {
	p.optDef = optDef
}

func (p *PokemonDto) OptSta() float64 {
	return p.optSta
}

func (p *PokemonDto) SetOptSta(optSta float64) {
	p.optSta = optSta
}
