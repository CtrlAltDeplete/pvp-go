package rawData

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	JSON_DIRECTORY = "../PvP-Go/rawData/json/"
	BASE_TYPES     = []string{"normal", "fire", "fighting", "water", "flying", "grass", "poison", "electric", "ground",
		"psychic", "rock", "ice", "bug", "dragon", "ghost", "dark", "steel", "fairy"}
	GEN_MAP = map[string]int64{
		"Generation 1": 1,
		"Generation 2": 2,
		"Generation 3": 3,
		"Generation 4": 4,
		"Generation 5": 5,
		"Generation 6": 6,
		"Generation 7": 7,
		"Sun and Moon": 7,
	}
	PVP_INELIGIBLE_LEGENDARIES = []string{"Mewtwo", "Kyogre", "Groudon", "Rayquaza", "Deoxys (Normal Forme)",
		"Deoxys (Attack Forme)", "Dialga", "Palkia", "Heatran", "Giratina (Origin Forme)", "Giratina (Altered Forme)"}
	PVP_ELIGIBLE_LEGENDARIES = []string{"Articuno", "Zapdos", "Moltres", "Mew", "Raikou", "Entei", "Suicune", "Lugia",
		"Ho-Oh", "Celebi", "Regirock", "Regice", "Registeel", "Latias", "Latios", "Deoxys (Defense Forme)",
		"Deoxys (Speed Forme)", "Uxie", "Mesprit", "Azelf", "Cresselia", "Meltan", "Melmetal"}
)

func DownloadJsonFiles() {
	var (
		TYPE_CHART_JSON           = "https://gamepress.gg/sites/default/files/aggregatedjson/POGOTypeChart.json"
		POKEMON_AND_MOVESETS_JSON = "https://gamepress.gg/sites/default/files/aggregatedjson/pokemon-data-full-en-PoGO.json"
		MOVES_JSON                = "https://gamepress.gg/sites/default/files/aggregatedjson/move-data-full-PoGO.json"
		err                       error
	)

	if _, err = os.Stat(JSON_DIRECTORY); os.IsNotExist(err) {
		err = os.MkdirAll(JSON_DIRECTORY, os.ModePerm)
		daos.CheckError(err)
	}

	downloadFile(TYPE_CHART_JSON, JSON_DIRECTORY+"type_chart.json")
	downloadFile(POKEMON_AND_MOVESETS_JSON, JSON_DIRECTORY+"pokemon_and_movesets.json")
	downloadFile(MOVES_JSON, JSON_DIRECTORY+"moves.json")
}

func downloadFile(url, path string) {
	var (
		resp *http.Response
		out  *os.File
		e    error
	)
	resp, e = http.Get(url)
	daos.CheckError(e)
	out, e = os.Create(path)
	daos.CheckError(e)
	_, e = io.Copy(out, resp.Body)
	daos.CheckError(e)
	daos.CheckError(resp.Body.Close())
	daos.CheckError(out.Close())
}

type typeChartDto struct {
	Name               string `json:"name"`
	FieldTypeAdvantage string `json:"field_type_advantage"`
	types              []string
	multipliers        map[string]float64
}

func FillTypeChartDto(dto *typeChartDto) {
	dto.types = strings.Split(strings.ToLower(dto.Name), "/")
	dto.multipliers = map[string]float64{}
	for _, t := range BASE_TYPES {
		dto.multipliers[t] = 1
	}
	r := regexp.MustCompile(`.*/([a-z]*)\.gif.*\n.*>((?:\d|\.)*)%</span> damage`)
	matches := r.FindAllStringSubmatch(dto.FieldTypeAdvantage, -1)
	for _, m := range matches {
		for i := 1; i < len(m); i += 2 {
			f, _ := strconv.ParseFloat(m[i+1], 64)
			dto.multipliers[m[i]] = f / 100.0
		}
	}
}

func ParseAllTypeCharts() {
	contents, err := ioutil.ReadFile(JSON_DIRECTORY + "type_chart.json")
	daos.CheckError(err)

	var typeChartDtos []typeChartDto
	err = json.Unmarshal(contents, &typeChartDtos)
	daos.CheckError(err)

	baseTypes := []*dtos.TypeDto{}
	for _, t := range BASE_TYPES {
		err, bt := daos.TYPES_DAO.Upsert(t, nil)
		daos.CheckError(err)
		baseTypes = append(baseTypes, bt)
	}

	var receivingType, actingType *dtos.TypeDto
	for _, dto := range typeChartDtos {
		FillTypeChartDto(&dto)
		if len(dto.types) == 1 {
			err, receivingType = daos.TYPES_DAO.Upsert(dto.types[0], nil)
		} else {
			err, receivingType = daos.TYPES_DAO.Upsert(dto.types[1], dto.types[1])
		}
		daos.CheckError(err)
		for _, actingType = range baseTypes {
			err, _ = daos.TYPE_MULTIPLIER_DAO.Upsert(receivingType.Id(), actingType.Id(), dto.multipliers[actingType.DisplayName()])
		}
	}
}

type moveDto struct {
	Title                  string `json:"title"`
	Power                  string `json:"power"`
	Cooldown               string `json:"cooldown"`
	MoveType               string `json:"move_type"` // TypeDto (toLower)
	Nid                    string `json:"nid"`
	EnergyGain             string `json:"energy_gain"`
	EnergyCost             string `json:"energy_cost"`
	DodgeWindow            string `json:"dodge_window"`
	MoveCategory           string `json:"move_category"` // "Fast MoveDto", "Charge MoveDto", ""
	DamageWindow           string `json:"damage_window"`
	PvpChargeEnergy        string `json:"pvp_charge_energy"`
	PvpChargeDamage        string `json:"pvp_charge_damage"`
	PvpFastDurationSeconds string `json:"pvp_fast_duration_seconds"`
	PvpFastDuration        string `json:"pvp_fast_duration"` // ++
	PvpFastEnergy          string `json:"pvp_fast_energy"`
	PvpFastPower           string `json:"pvp_fast_power"`
	TitleLinked            string `json:"title_linked"`
	Probability            string `json:"probability"` // "", 0.10, 1.00, 0.50, 0.30, 0.13
	StageDelta             string `json:"stage_delta"` // "", "2", "-1", "1", "-2"
	Stat                   string `json:"stat"`        // "", "Atk, Def", "Atk", "Def"
	Subject                string `json:"subject"`     // "", "Self", "Opponent"
	name                   string
	typeId                 int64
	power                  int64
	turns                  int64
	energy                 int64
	probability            sql.NullFloat64
	stageDelta             sql.NullInt64
	stats                  sql.NullString
	target                 sql.NullString
}

func FillMoveDto(dto *moveDto) {
	dto.name = dto.Title
	err, pokemonType := daos.TYPES_DAO.FindSingleByType(strings.ToLower(dto.MoveType))
	if err != nil {
		log.Printf("No type %s for move %s.\n", strings.ToLower(dto.MoveType), dto.Title)
		return
	}
	dto.typeId = pokemonType.Id()

	if dto.MoveCategory == "Fast MoveDto" {
		dto.power, err = strconv.ParseInt(dto.PvpFastPower, 10, 64)
		daos.CheckError(err)

		dto.turns, err = strconv.ParseInt(dto.PvpFastDuration, 10, 64)
		daos.CheckError(err)
		dto.turns += 1

		dto.energy, err = strconv.ParseInt(dto.PvpFastEnergy, 10, 64)
		daos.CheckError(err)

		dto.probability.Valid = false
		dto.stageDelta.Valid = false
		dto.stats.Valid = false
		dto.target.Valid = false
	} else if dto.MoveCategory == "Charge MoveDto" {
		dto.power, err = strconv.ParseInt(dto.PvpChargeDamage, 10, 64)
		daos.CheckError(err)

		dto.turns = 1

		dto.energy, err = strconv.ParseInt(dto.PvpChargeEnergy, 10, 64)
		daos.CheckError(err)

		dto.probability.Valid = dto.Probability != ""
		if dto.probability.Valid {
			dto.probability.Float64, err = strconv.ParseFloat(dto.Probability, 64)
			daos.CheckError(err)
		}

		dto.stageDelta.Valid = dto.StageDelta != ""
		if dto.stageDelta.Valid {
			dto.stageDelta.Int64, err = strconv.ParseInt(dto.StageDelta, 10, 64)
			daos.CheckError(err)
		}

		dto.stats.String = dto.Stat
		dto.stats.Valid = dto.stats.String != ""

		dto.target.String = dto.Subject
		dto.target.Valid = dto.target.String != ""
	} else {
		log.Printf("Incomplete data for move %s.\n", dto.Title)
		return
	}
}

func ParseAllMoves() {
	contents, e := ioutil.ReadFile(JSON_DIRECTORY + "moves.json")
	daos.CheckError(e)

	var moveDtos []moveDto
	e = json.Unmarshal(contents, &moveDtos)
	daos.CheckError(e)

	for _, dto := range moveDtos {
		FillMoveDto(&dto)
		err, _ := daos.MOVES_DAO.Upsert(dto.name, dto.typeId, dto.power, dto.turns, dto.energy, dto.probability,
			dto.stageDelta, dto.stats, dto.target)
		daos.CheckError(err)
	}
}

type pokemonAndMovesetsDto struct {
	AlternateForm           string `json:"alternate_form"`
	Number                  string `json:"number"`
	Sta                     string `json:"sta"`
	Atk                     string `json:"atk"`
	Def                     string `json:"def"`
	Cp                      string `json:"cp"`
	Rating                  string `json:"rating"`
	FieldPokemonGeneration  string `json:"field_pokemon_generation"`
	FieldPokemonType        string `json:"field_pokemon_type"`
	Title1                  string `json:"title_1"`
	Uri                     string `json:"uri"`
	FieldPrimaryMoves       string `json:"field_primary_moves"`
	FieldSecondaryMoves     string `json:"field_secondary_moves"`
	Nid                     string `json:"nid"`
	FieldLegacyChargeMoves  string `json:"field_legacy_charge_moves"`
	FieldLegacyQuickMoves   string `json:"field_legacy_quick_moves"`
	FieldEvolutions         string `json:"field_evolutions"`
	CatchRate               string `json:"catch_rate"`
	FieldFleeRate           string `json:"field_flee_rate"`
	QuickExclusiveMoves     string `json:"quick_exclusive_moves"`
	ChargeExclusiveMoves    string `json:"charge_exclusive_moves"`
	SecondaryChargeCost     string `json:"second_charge_cost"`
	PokemonImageSmall       string `json:"pokemon_image_small"`
	Title                   string `json:"title"`
	Path                    string `json:"path"`
	Lvl20                   string `json:"lvl20"`
	Lvl30                   string `json:"lvl30"`
	Lvl40                   string `json:"lvl40"`
	gen                     int64
	name                    string
	typeId                  int64
	atk                     float64
	def                     float64
	sta                     float64
	dateAdd                 string
	isLegendary             bool
	isPvpEligible           bool
	fastMoveIdAndIsLegacy   map[int64]bool
	chargeMoveIdAndIsLegacy map[int64]bool
}

func FillPokemonAndMovesetDto(dto *pokemonAndMovesetsDto) error {
	dto.gen = GEN_MAP[dto.FieldPokemonGeneration]
	dto.name = dto.Title1
	types := strings.Split(strings.ToLower(dto.FieldPokemonType), ", ")

	var err error
	var pokemonType *dtos.TypeDto
	err, pokemonType = daos.TYPES_DAO.FindSingleByTypes(types)
	if err == daos.NO_ROWS {
		log.Printf("Cannot find TypeDto for %v for %s.\n", types, dto.name)
		return err
	} else if err != nil {
		return err
	}
	dto.typeId = pokemonType.Id()

	dto.atk, err = strconv.ParseFloat(dto.Atk, 64)
	daos.CheckError(err)

	dto.def, err = strconv.ParseFloat(dto.Def, 64)
	daos.CheckError(err)

	dto.sta, err = strconv.ParseFloat(dto.Sta, 64)
	daos.CheckError(err)

	r := regexp.MustCompile(`.*/pokemongo/files/(\d*)-(\d*).*`)
	matches := r.FindAllStringSubmatch(dto.Uri, -1)[0]
	dto.dateAdd = fmt.Sprintf("%s-%s-01", matches[1], matches[2])

	dto.isLegendary = false
	dto.isPvpEligible = true
	for _, name := range PVP_ELIGIBLE_LEGENDARIES {
		if dto.name == name {
			dto.isLegendary = true
		}
	}
	for _, name := range PVP_INELIGIBLE_LEGENDARIES {
		if dto.name == name {
			dto.isLegendary = true
			dto.isPvpEligible = false
		}
	}

	dto.fastMoveIdAndIsLegacy = map[int64]bool{}
	var move *dtos.MoveDto
	dto.chargeMoveIdAndIsLegacy = map[int64]bool{}

	for _, fastName := range strings.Split(dto.FieldPrimaryMoves, ", ") {
		if fastName == "" {
			continue
		}
		err, move = daos.MOVES_DAO.FindByName(fastName)
		if err == daos.NO_ROWS {
			log.Printf("Could not find move %s for pokemon %s.\n", fastName, dto.name)
		} else if err != nil {
			log.Print(err.Error())
		} else {
			dto.fastMoveIdAndIsLegacy[move.Id()] = false
		}
	}

	for _, chargeName := range strings.Split(dto.FieldSecondaryMoves, ", ") {
		if chargeName == "" {
			continue
		}
		err, move = daos.MOVES_DAO.FindByName(chargeName)
		if err == daos.NO_ROWS {
			log.Printf("Could not find move %s for pokemon %s.\n", chargeName, dto.name)
		} else if err != nil {
			log.Print(err.Error())
		} else {
			dto.chargeMoveIdAndIsLegacy[move.Id()] = false
		}
	}

	for _, fastName := range strings.Split(dto.FieldLegacyQuickMoves, ", ") {
		if fastName == "" {
			continue
		}
		err, move = daos.MOVES_DAO.FindByName(fastName)
		if err == daos.NO_ROWS {
			log.Printf("Could not find move %s for pokemon %s.\n", fastName, dto.name)
		} else if err != nil {
			log.Print(err.Error())
		} else {
			dto.fastMoveIdAndIsLegacy[move.Id()] = true
		}
	}

	for _, chargeName := range strings.Split(dto.FieldLegacyChargeMoves, ", ") {
		if chargeName == "" {
			continue
		}
		err, move = daos.MOVES_DAO.FindByName(chargeName)
		if err == daos.NO_ROWS {
			log.Printf("Could not find move %s for pokemon %s.\n", chargeName, dto.name)
		} else if err != nil {
			log.Print(err.Error())
		} else {
			dto.chargeMoveIdAndIsLegacy[move.Id()] = true
		}
	}

	for _, fastName := range strings.Split(dto.QuickExclusiveMoves, ", ") {
		if fastName == "" {
			continue
		}
		err, move = daos.MOVES_DAO.FindByName(fastName)
		if err == daos.NO_ROWS {
			log.Printf("Could not find move %s for pokemon %s.\n", fastName, dto.name)
		} else if err != nil {
			log.Print(err.Error())
		} else {
			dto.fastMoveIdAndIsLegacy[move.Id()] = true
		}
	}

	for _, chargeName := range strings.Split(dto.ChargeExclusiveMoves, ", ") {
		if chargeName == "" {
			continue
		}
		err, move = daos.MOVES_DAO.FindByName(chargeName)
		if err == daos.NO_ROWS {
			log.Printf("Could not find move %s for pokemon %s.\n", chargeName, dto.name)
		} else if err != nil {
			log.Print(err.Error())
		} else {
			dto.chargeMoveIdAndIsLegacy[move.Id()] = true
		}
	}

	if len(dto.chargeMoveIdAndIsLegacy) == 0 || len(dto.fastMoveIdAndIsLegacy) == 0 {
		log.Printf("Could not find moves for pokemon %s.\n", dto.name)
		return errors.New(fmt.Sprintf("Could not find moves for pokemon %s.\n", dto.name))
	}

	return nil
}

func ParseAllPokemonAndMovesets() {
	contents, e := ioutil.ReadFile(JSON_DIRECTORY + "pokemon_and_movesets.json")
	daos.CheckError(e)

	var pokemonAndMovesetsDtos []pokemonAndMovesetsDto
	e = json.Unmarshal(contents, &pokemonAndMovesetsDtos)
	daos.CheckError(e)

	var pokemon *dtos.PokemonDto

	for i, dto := range pokemonAndMovesetsDtos {
		err := FillPokemonAndMovesetDto(&dto)
		if err != nil {
			log.Printf("[%d/%d] %f%% complete.\n", i, len(pokemonAndMovesetsDtos), (100.0 * float64(i) / float64(len(pokemonAndMovesetsDtos))))
			continue
		}
		err, pokemon = daos.POKEMON_DAO.Upsert(dto.gen, dto.name, dto.typeId, dto.atk, dto.def, dto.sta, dto.dateAdd,
			dto.isLegendary, dto.isPvpEligible, 0, 0, 0, 0)
		daos.CheckError(err)

		for moveId, isMoveLegacy := range dto.fastMoveIdAndIsLegacy {
			err, _ = daos.POKEMON_HAS_MOVE_DAO.Upsert(pokemon.Id(), moveId, isMoveLegacy)
			daos.CheckError(err)
		}
		for moveId, isMoveLegacy := range dto.chargeMoveIdAndIsLegacy {
			err, _ = daos.POKEMON_HAS_MOVE_DAO.Upsert(pokemon.Id(), moveId, isMoveLegacy)
			daos.CheckError(err)
		}
		log.Printf("[%d/%d] %f%% complete.\n", i, len(pokemonAndMovesetsDtos), (100.0 * float64(i) / float64(len(pokemonAndMovesetsDtos))))
	}
}
