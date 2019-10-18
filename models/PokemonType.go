package models

import (
	"database/sql"
	"log"
)

type PokemonType struct {
	id          int64
	firstType   string
	secondType  sql.NullString
	displayName string
}

func (p *PokemonType) Id() int64 {
	return p.id
}

func (p *PokemonType) SetId(id int64) {
	p.id = id
}

func (p *PokemonType) FirstType() string {
	return p.firstType
}

func (p *PokemonType) SetFirstType(firstType string) {
	p.firstType = firstType
	p.updateDisplayName()
}

func (p *PokemonType) IsSecondTypeNull() bool {
	return !p.secondType.Valid
}

func (p *PokemonType) SecondType() string {
	return p.secondType.String
}

func (p *PokemonType) SetSecondType(secondType interface{}) {
	switch st := secondType.(type) {
	case string:
		p.secondType.Valid = true
		p.secondType.String = st
	case nil:
		p.secondType.Valid = false
		p.secondType.String = ""
	case sql.NullString:
		if st.Valid {
			p.SetSecondType(st.String)
		} else {
			p.SetSecondType(nil)
		}
	default:
		log.Fatalf("Unknown type %T.", st)
	}
	p.updateDisplayName()
}

func (p *PokemonType) DisplayName() string {
	return p.displayName
}

func (p *PokemonType) SetDisplayName(displayName string) {
	p.displayName = displayName
}

func (p *PokemonType) updateDisplayName() {
	if p.IsSecondTypeNull() {
		p.SetDisplayName(p.FirstType())
	} else {
		p.SetDisplayName(p.FirstType() + "/" + p.SecondType())
	}
}
