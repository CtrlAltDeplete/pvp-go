package main

import (
	"fmt"
	"testing"
)

func TestBattle_BulbasaurVsSquirtle(t *testing.T) {
	var (
		bulbasaur, squirtle                                                        Pokemon
		bulbasaurDto, squirtleDto                                                  *api.PokemonDto
		vineWhipDto, powerWhipDto, sludgeBombDto, bubbleDto, aquaJetDto, returnDto *api.MoveDto
		battle                                                                     Battle
		startingShields                                                            int64 = 1
		bulbasaurScore, squirtleScore                                              int64
		err                                                                        error
	)
	err, bulbasaurDto = api.POKEMON_DAO.FindByName("Bulbasaur")
	api.CheckError(err)
	err, vineWhipDto = api.MOVES_DAO.FindByName("Vine Whip")
	api.CheckError(err)
	err, powerWhipDto = api.MOVES_DAO.FindByName("Power Whip")
	api.CheckError(err)
	err, sludgeBombDto = api.MOVES_DAO.FindByName("Sludge Bomb")
	api.CheckError(err)

	err, squirtleDto = api.POKEMON_DAO.FindByName("Squirtle")
	api.CheckError(err)
	err, bubbleDto = api.MOVES_DAO.FindByName("Bubble")
	api.CheckError(err)
	err, aquaJetDto = api.MOVES_DAO.FindByName("Aqua Jet")
	api.CheckError(err)
	err, returnDto = api.MOVES_DAO.FindByName("Return")
	api.CheckError(err)

	bulbasaur = *NewPokemon(*bulbasaurDto, *vineWhipDto, []api.MoveDto{*powerWhipDto, *sludgeBombDto})
	squirtle = *NewPokemon(*squirtleDto, *bubbleDto, []api.MoveDto{*aquaJetDto, *returnDto})

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
		medichamDto, skarmoryDto                               *api.PokemonDto
		counterDto, powerUpPunchDto, airSlashDto, skyAttackDto *api.MoveDto
		battle                                                 Battle
		startingShields                                        int64 = 2
		medichamScore, skarmoryScore                           int64
		err                                                    error
	)
	err, medichamDto = api.POKEMON_DAO.FindByName("Medicham")
	api.CheckError(err)
	err, counterDto = api.MOVES_DAO.FindByName("Counter")
	api.CheckError(err)
	err, powerUpPunchDto = api.MOVES_DAO.FindByName("Power-Up Punch")
	api.CheckError(err)

	err, skarmoryDto = api.POKEMON_DAO.FindByName("Skarmory")
	api.CheckError(err)
	err, airSlashDto = api.MOVES_DAO.FindByName("Air Slash")
	api.CheckError(err)
	err, skyAttackDto = api.MOVES_DAO.FindByName("Sky Attack")
	api.CheckError(err)

	medicham = *NewPokemon(*medichamDto, *counterDto, []api.MoveDto{*powerUpPunchDto})
	skarmory = *NewPokemon(*skarmoryDto, *airSlashDto, []api.MoveDto{*skyAttackDto})

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
		quagsireDto, omastarDto                              *api.PokemonDto
		mudShotDto, acidSprayDto, rockThrowDto, hydroPumpDto *api.MoveDto
		battle                                               Battle
		startingShields                                      int64 = 1
		quagsireScore, omastarScore                          int64
		err                                                  error
	)
	err, quagsireDto = api.POKEMON_DAO.FindByName("Quagsire")
	api.CheckError(err)
	err, mudShotDto = api.MOVES_DAO.FindByName("Mud Shot")
	api.CheckError(err)
	err, acidSprayDto = api.MOVES_DAO.FindByName("Acid Spray")
	api.CheckError(err)

	err, omastarDto = api.POKEMON_DAO.FindByName("Omastar")
	api.CheckError(err)
	err, rockThrowDto = api.MOVES_DAO.FindByName("Rock Throw")
	api.CheckError(err)
	err, hydroPumpDto = api.MOVES_DAO.FindByName("Hydro Pump")
	api.CheckError(err)

	quagsire = *NewPokemon(*quagsireDto, *mudShotDto, []api.MoveDto{*acidSprayDto})
	omastar = *NewPokemon(*omastarDto, *rockThrowDto, []api.MoveDto{*hydroPumpDto})

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
