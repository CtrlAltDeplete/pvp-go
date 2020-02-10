package main

import (
	"database/sql"
	"log"
)

type TypeDto struct {
	id          int64
	firstType   string
	secondType  sql.NullString
	displayName string
}

func (p *TypeDto) Id() int64 {
	return p.id
}

func (p *TypeDto) SetId(id int64) {
	p.id = id
}

func (p *TypeDto) FirstType() string {
	return p.firstType
}

func (p *TypeDto) SetFirstType(firstType string) {
	p.firstType = firstType
	p.updateDisplayName()
}

func (p *TypeDto) IsSecondTypeNull() bool {
	return !p.secondType.Valid
}

func (p *TypeDto) SecondType() string {
	return p.secondType.String
}

func (p *TypeDto) SecondTypeNullable() sql.NullString {
	return p.secondType
}

func (p *TypeDto) SetSecondType(secondType interface{}) {
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

func (p *TypeDto) DisplayName() string {
	return p.displayName
}

func (p *TypeDto) SetDisplayName(displayName string) {
	p.displayName = displayName
}

func (p *TypeDto) updateDisplayName() {
	if p.IsSecondTypeNull() {
		p.SetDisplayName(p.FirstType())
	} else {
		p.SetDisplayName(p.FirstType() + "/" + p.SecondType())
	}
}
