package models

import (
	"PvP-Go/db/daos"
	"PvP-Go/db/dtos"
	"fmt"
	"testing"
)

func TestBattle_BulbasaurVsSquirtle(t *testing.T) {
	var (
		bulbasaur, squirtle                                                        Pokemon
		bulbasaurDto, squirtleDto                                                  *dtos.PokemonDto
		vineWhipDto, powerWhipDto, sludgeBombDto, bubbleDto, aquaJetDto, returnDto *dtos.MoveDto
		battle                                                                     Battle
		startingShields                                                            int64 = 1
		bulbasaurScore, squirtleScore                                              int64
		err                                                                        error
	)
	err, bulbasaurDto = daos.POKEMON_DAO.FindByName("Bulbasaur")
	daos.CheckError(err)
	err, vineWhipDto = daos.MOVES_DAO.FindByName("Vine Whip")
	daos.CheckError(err)
	err, powerWhipDto = daos.MOVES_DAO.FindByName("Power Whip")
	daos.CheckError(err)
	err, sludgeBombDto = daos.MOVES_DAO.FindByName("Sludge Bomb")
	daos.CheckError(err)

	err, squirtleDto = daos.POKEMON_DAO.FindByName("Squirtle")
	daos.CheckError(err)
	err, bubbleDto = daos.MOVES_DAO.FindByName("Bubble")
	daos.CheckError(err)
	err, aquaJetDto = daos.MOVES_DAO.FindByName("Aqua Jet")
	daos.CheckError(err)
	err, returnDto = daos.MOVES_DAO.FindByName("Return")
	daos.CheckError(err)

	bulbasaur = *NewPokemon(*bulbasaurDto, *vineWhipDto, []dtos.MoveDto{*powerWhipDto, *sludgeBombDto})
	squirtle = *NewPokemon(*squirtleDto, *bubbleDto, []dtos.MoveDto{*aquaJetDto, *returnDto})

	battle = *NewBattle([]Pokemon{bulbasaur, squirtle}, []int64{startingShields, startingShields})
	bulbasaurScore, squirtleScore = battle.Simulate()
	if battle.pokemon[0].hp != 51 || battle.pokemon[1].hp != 0 {
		fmt.Printf("Expected Bulbasaur %d hp, Squirtle %d hp; got Bulbasaur %d hp, Squirtle %d hp\n",
			51, 0, int(battle.pokemon[0].hp), int(battle.pokemon[1].hp))
		t.Fail()
	} else {
		fmt.Printf("Score: Bulbasaur %d vs Squirtle %d\n", bulbasaurScore, squirtleScore)
	}
}

func TestBattle_MedichamVsSkarmory(t *testing.T) {
	var (
		medicham, skarmory                                     Pokemon
		medichamDto, skarmoryDto                               *dtos.PokemonDto
		counterDto, powerUpPunchDto, airSlashDto, skyAttackDto *dtos.MoveDto
		battle                                                 Battle
		startingShields                                        int64 = 2
		medichamScore, skarmoryScore                           int64
		err                                                    error
	)
	err, medichamDto = daos.POKEMON_DAO.FindByName("Medicham")
	daos.CheckError(err)
	err, counterDto = daos.MOVES_DAO.FindByName("Counter")
	daos.CheckError(err)
	err, powerUpPunchDto = daos.MOVES_DAO.FindByName("Power-Up Punch")
	daos.CheckError(err)

	err, skarmoryDto = daos.POKEMON_DAO.FindByName("Skarmory")
	daos.CheckError(err)
	err, airSlashDto = daos.MOVES_DAO.FindByName("Air Slash")
	daos.CheckError(err)
	err, skyAttackDto = daos.MOVES_DAO.FindByName("Sky Attack")
	daos.CheckError(err)

	medicham = *NewPokemon(*medichamDto, *counterDto, []dtos.MoveDto{*powerUpPunchDto})
	skarmory = *NewPokemon(*skarmoryDto, *airSlashDto, []dtos.MoveDto{*skyAttackDto})

	battle = *NewBattle([]Pokemon{medicham, skarmory}, []int64{startingShields, startingShields})
	medichamScore, skarmoryScore = battle.Simulate()
	if battle.pokemon[0].hp != 6 || battle.pokemon[1].hp != 0 {
		fmt.Printf("Expected Medicham %d hp, Skarmory %d hp; got Medicham %d hp, Skarmory %d hp\n",
			6, 0, int(battle.pokemon[0].hp), int(battle.pokemon[1].hp))
		t.Fail()
	} else {
		fmt.Printf("Score: Medicham %d vs Skarmory %d\n", medichamScore, skarmoryScore)
	}
}

func TestBattle_QuagsireVsOmastar(t *testing.T) {
	var (
		quagsire, omastar                                    Pokemon
		quagsireDto, omastarDto                              *dtos.PokemonDto
		mudShotDto, acidSprayDto, rockThrowDto, hydroPumpDto *dtos.MoveDto
		battle                                               Battle
		startingShields                                      int64 = 1
		quagsireScore, omastarScore                          int64
		err                                                  error
	)
	err, quagsireDto = daos.POKEMON_DAO.FindByName("Quagsire")
	daos.CheckError(err)
	err, mudShotDto = daos.MOVES_DAO.FindByName("Mud Shot")
	daos.CheckError(err)
	err, acidSprayDto = daos.MOVES_DAO.FindByName("Acid Spray")
	daos.CheckError(err)

	err, omastarDto = daos.POKEMON_DAO.FindByName("Omastar")
	daos.CheckError(err)
	err, rockThrowDto = daos.MOVES_DAO.FindByName("Rock Throw")
	daos.CheckError(err)
	err, hydroPumpDto = daos.MOVES_DAO.FindByName("Hydro Pump")
	daos.CheckError(err)

	quagsire = *NewPokemon(*quagsireDto, *mudShotDto, []dtos.MoveDto{*acidSprayDto})
	omastar = *NewPokemon(*omastarDto, *rockThrowDto, []dtos.MoveDto{*hydroPumpDto})

	battle = *NewBattle([]Pokemon{quagsire, omastar}, []int64{startingShields, startingShields})
	quagsireScore, omastarScore = battle.Simulate()
	if battle.pokemon[0].hp != 69 || battle.pokemon[1].hp != 0 {
		fmt.Printf("Expected Quagsire %d hp, Omastar %d hp; got Quagsire %d hp, Omastar %d hp\n",
			64, 0, int(battle.pokemon[0].hp), int(battle.pokemon[1].hp))
		t.Fail()
	} else {
		fmt.Printf("Score: Quagsire %d vs Omastar %d\n", quagsireScore, omastarScore)
	}
}
