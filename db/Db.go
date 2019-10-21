package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

var (
	LIVE                 = NewDB("user", "password", "host")
	NO_ROWS              = errors.New("No rows found.")
	MULTIPLE_ROWS        = errors.New("Multiple rows found.")
	TYPE_MISMATCH        = errors.New("Incorrect parameter type for creation.")
	MOVES_DAO            = MovesDao{}
	POKEMON_DAO          = PokemonDao{}
	POKEMON_HAS_MOVE_DAO = PokemonHasMoveDao{}
	TYPE_MULTIPLIER_DAO  = TypeMultipliersDao{}
	TYPES_DAO            = TypesDao{}
)

func CheckError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func NewDB(user, password, address string) *sql.DB {
	db, e := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/pvpgo", user, password, address))
	CheckError(e)
	return db
}
