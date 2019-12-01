package main

import (
	"PvP-Go/models"
	"fmt"
)

func main() {
	fmt.Print("Starting Sinister")
	models.NewCup("sinister", models.SINISTER_SQL).CalculateOtherTables()
	fmt.Print("Starting Ferocious")
	models.NewCup("ferocious", models.FEROCIOUS_SQL).CalculateOtherTables()
	fmt.Print("Starting Timeless")
	models.NewCup("timeless", models.TIMELESS_SQL).CalculateOtherTables()
}
