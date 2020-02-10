package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

var (
	LIVE = NewDB(os.Getenv("PVP_GO_DB_USER"),
		os.Getenv("PVP_GO_DB_PASSWORD"),
		os.Getenv("PVP_GO_DB_ENDPOINT"))
	NO_ROWS              = errors.New("No rows found.")
	MULTIPLE_ROWS        = errors.New("Multiple rows found.")
	BAD_PARAMS           = errors.New("Bad parameters for DAO.")
	API_DAO              = ApiDao{}
	BATTLE_SIMS_DAO      = BattleSimulationsDao{}
	CP_DAO               = CpDao{}
	MATCH_UPS_DAO        = MatchUpsDao{}
	MOVES_DAO            = MovesDao{}
	MOVE_SETS_DAO        = MoveSetDao{}
	POKEMON_DAO          = PokemonDao{}
	POKEMON_HAS_MOVE_DAO = PokemonHasMoveDao{}
	RANKINGS_DAO         = RankingsDao{}
	TEAM_RANKINGS_DAO    = TeamRankingsDao{}
	TYPE_MULTIPLIER_DAO  = TypeMultipliersDao{}
	TYPES_DAO            = TypesDao{}
)

func CheckError(e error) {
	if e != nil {
		log.Panic(e)
	}
}

func NewDB(user, password, address string) *sql.DB {
	db, e := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/pvpgo", user, password, address))
	CheckError(e)
	return db
}
