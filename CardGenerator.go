package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"os"
)

func BuildCard(moveSet MoveSetDto, cup string, cupTypes []string) ([]byte, error) {
	_, pokemonDto := POKEMON_DAO.FindById(moveSet.PokemonId())
	_, fastMoveDto := MOVES_DAO.FindById(moveSet.FastMoveId())
	_, primaryChargeDto := MOVES_DAO.FindById(moveSet.PrimaryChargeMoveId())
	chargeMoves := []MoveDto{*primaryChargeDto}
	if moveSet.SecondaryChargeMoveId() != nil {
		_, secondaryChargeDto := MOVES_DAO.FindById(*moveSet.SecondaryChargeMoveId())
		chargeMoves = append(chargeMoves, *secondaryChargeDto)
	}
	_, typesDto := TYPES_DAO.FindById(pokemonDto.TypeId())

	pokemon := NewPokemon(*pokemonDto, *fastMoveDto, chargeMoves)

	var (
		fontfile = "./rawData/imgs/GillSansStd-Bold.ttf"
	)

	img := image.NewRGBA(image.Rect(0, 0, 500, 700))
	textColor, err := applyTypeProperties(img, cupTypes, *typesDto)
	if err != nil {
		return nil, err
	}
	err = applyCupProperties(img, cup)
	if err != nil {
		return nil, err
	}

	fontBytes, err := ioutil.ReadFile(fontfile)
	if err != nil {
		log.Fatal(err)
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Fatal(err)
	}

	err = applyPokemon(img, *pokemonDto, *f, textColor)
	if err != nil {
		return nil, err
	}

	err = applyFastMove(img, *pokemon, *f, textColor)
	if err != nil {
		return nil, err
	}

	err = applyPrimaryCharge(img, *pokemon, *f, textColor)
	if err != nil {
		return nil, err
	}

	if moveSet.SecondaryChargeMoveId() != nil {
		err = applySecondaryCharge(img, *pokemon, *f, textColor)
		if err != nil {
			return nil, err
		}
	}

	err = applyMatchups(img, moveSet.Id(), cup)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	err = png.Encode(buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func applyTypeProperties(img *image.RGBA, cupTypes []string, pokemonTypes TypeDto) (*image.Uniform, error) {
	var (
		cardFile       *os.File
		cardBackground image.Image
		textColor      = image.Black
		err            error
	)
	if pokemonTypes.IsSecondTypeNull() {
		cardFile, err = os.Open(fmt.Sprintf("./rawData/imgs/cardBg/%s.png", pokemonTypes.FirstType()))
		if err != nil {
			return nil, err
		}
		cardBackground, err = png.Decode(cardFile)
		if err != nil {
			return nil, err
		}
		for _, darkColor := range []string{"dark", "dragon", "fighting", "ghost", "steel"} {
			if darkColor == pokemonTypes.FirstType() {
				textColor = image.White
			}
		}
		draw.Draw(img, cardBackground.Bounds(), cardBackground, image.Pt(0, 0), draw.Src)
	} else {
		var eligibleTypes []string
		for _, t := range cupTypes {
			if t == pokemonTypes.FirstType() || t == pokemonTypes.SecondType() {
				eligibleTypes = append(eligibleTypes, t)
			}
		}
		var primaryType, secondaryType string
		if len(eligibleTypes) == 1 {
			primaryType = eligibleTypes[0]
			if pokemonTypes.FirstType() == primaryType {
				secondaryType = pokemonTypes.SecondType()
			} else {
				secondaryType = pokemonTypes.FirstType()
			}
		} else {
			if rand.Int()%2 == 0 {
				primaryType = pokemonTypes.FirstType()
				secondaryType = pokemonTypes.SecondType()
			} else {
				primaryType = pokemonTypes.SecondType()
				secondaryType = pokemonTypes.FirstType()
			}
		}

		cardFile, err = os.Open(fmt.Sprintf("./rawData/imgs/cardBg/%s.png", primaryType))
		if err != nil {
			return nil, err
		}
		cardBackground, err = png.Decode(cardFile)
		if err != nil {
			return nil, err
		}
		for _, darkColor := range []string{"dark", "dragon", "fighting", "ghost", "steel"} {
			if darkColor == primaryType {
				textColor = image.White
			}
		}
		draw.Draw(img, cardBackground.Bounds(), cardBackground, image.Pt(0, 0), draw.Src)

		var (
			typeIconFile *os.File
			typeIconImg  image.Image
		)

		typeIconFile, err = os.Open(fmt.Sprintf("./rawData/imgs/typeIcons/%s-icon.png", secondaryType))
		if err != nil {
			return nil, err
		}
		typeIconImg, err = png.Decode(typeIconFile)
		if err != nil {
			return nil, err
		}
		draw.Draw(img, image.Rect(370, 25, 415, 65), typeIconImg, image.Pt(0, 0), draw.Over)
	}
	return textColor, nil
}

func applyCupProperties(img *image.RGBA, cup string) error {
	var (
		cupIconFile *os.File
		cupIconImg  image.Image
		cupBgFile   *os.File
		cupBgImg    image.Image
		borderFile  *os.File
		borderImg   image.Image
		err         error
	)
	cupIconFile, err = os.Open(fmt.Sprintf("./rawData/imgs/cupIcons/%s-icon.png", cup))
	if err != nil {
		return err
	}
	cupIconImg, err = png.Decode(cupIconFile)
	if err != nil {
		return err
	}
	draw.Draw(img, image.Rect(40, 25, 80, 65), cupIconImg, image.Pt(0, 0), draw.Over)

	cupBgFile, err = os.Open(fmt.Sprintf("./rawData/imgs/cupBg/%s.png", cup))
	if err != nil {
		return err
	}
	cupBgImg, err = png.Decode(cupBgFile)
	if err != nil {
		return err
	}
	draw.Draw(img, image.Rect(40, 70, 460, 335), cupBgImg, image.Pt(0, 0), draw.Over)

	borderFile, err = os.Open("./rawData/imgs/cardBg/border.png")
	if err != nil {
		return err
	}
	borderImg, err = png.Decode(borderFile)
	if err != nil {
		return err
	}
	draw.Draw(img, img.Bounds(), borderImg, image.Pt(0, 0), draw.Over)
	return nil
}

func applyPokemon(img *image.RGBA, pokemon PokemonDto, ttf truetype.Font, textColor *image.Uniform) error {
	var (
		pokemonFile                                   *os.File
		pokemonImg                                    image.Image
		width, height                                 int
		widthMultiplier, heightMultiplier, multiplier float64
		targetWidth                                   = 380.0
		targetHeight                                  = 225.0
		startX                                        = 60
		startY                                        = 90
		newWidth, newHeight                           int
		err                                           error
	)

	pokemonFile, err = os.Open(fmt.Sprintf("./rawData/imgs/pokemon/%d.png", pokemon.Id()))
	if err != nil {
		return err
	}
	pokemonImg, err = png.Decode(pokemonFile)
	if err != nil {
		return err
	}

	bounds := pokemonImg.Bounds()
	width = bounds.Size().X
	height = bounds.Size().Y
	widthMultiplier = targetWidth / float64(width)
	heightMultiplier = targetHeight / float64(height)
	multiplier = math.Min(widthMultiplier, heightMultiplier)
	newWidth = int(float64(width) * multiplier)
	newHeight = int(float64(height) * multiplier)
	startX += (int(targetWidth) - newWidth) / 2
	startY += (int(targetHeight) - newHeight) / 2

	draw.BiLinear.Scale(img, image.Rect(startX, startY, startX+newWidth, startY+newHeight), pokemonImg, pokemonImg.Bounds(), draw.Over, nil)

	h := font.HintingNone
	drawer := font.Drawer{
		Dst: img,
		Src: textColor,
		Face: truetype.NewFace(&ttf, &truetype.Options{
			Size:    30.0,
			DPI:     75.0,
			Hinting: h,
		}),
	}

	drawer.Dot = fixed.Point26_6{
		X: fixed.I(90),
		Y: fixed.I(57),
	}
	drawer.DrawString(pokemon.Name())
	return nil
}

func applyFastMove(img *image.RGBA, pokemon Pokemon, ttf truetype.Font, textColor *image.Uniform) error {
	var (
		typeIconFile *os.File
		typeIconImg  image.Image
		move         = pokemon.FastMove()
		err          error
	)
	_, typeDto := TYPES_DAO.FindById(move.TypeId())

	typeIconFile, err = os.Open(fmt.Sprintf("./rawData/imgs/typeIcons/%s-icon.png", typeDto.DisplayName()))
	if err != nil {
		return err
	}
	typeIconImg, err = png.Decode(typeIconFile)
	if err != nil {
		return err
	}

	draw.Draw(img, image.Rect(45, 355, 85, 395), typeIconImg, image.Pt(0, 0), draw.Over)

	h := font.HintingNone
	nameDrawer := font.Drawer{
		Dst: img,
		Src: textColor,
		Face: truetype.NewFace(&ttf, &truetype.Options{
			Size:    24.0,
			DPI:     75.0,
			Hinting: h,
		}),
	}
	nameDrawer.Dot = fixed.Point26_6{
		X: fixed.I(90),
		Y: fixed.I(373),
	}
	nameDrawer.DrawString(move.Name())

	extraDrawer := font.Drawer{
		Dst: img,
		Src: textColor,
		Face: truetype.NewFace(&ttf, &truetype.Options{
			Size:    16.0,
			DPI:     75.0,
			Hinting: h,
		}),
	}
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(90),
		Y: fixed.I(391),
	}
	extraDrawer.DrawString(fmt.Sprintf("Turns: %d", move.CoolDown()))

	power := move.Power() * pokemon.GetStab(&move)
	powerLabel := fmt.Sprintf("Power: %.1f", power)
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(355+52) - extraDrawer.MeasureString(powerLabel)/2,
		Y: fixed.I(370),
	}
	extraDrawer.DrawString(powerLabel)

	energyLabel := fmt.Sprintf("Energy: %.0f", move.Energy())
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(355+52) - extraDrawer.MeasureString(energyLabel)/2,
		Y: fixed.I(391),
	}
	extraDrawer.DrawString(energyLabel)
	return nil
}

func applyPrimaryCharge(img *image.RGBA, pokemon Pokemon, ttf truetype.Font, textColor *image.Uniform) error {
	var (
		typeIconFile *os.File
		typeIconImg  image.Image
		move         = pokemon.ChargeMoves()[0]
		err          error
	)
	_, typeDto := TYPES_DAO.FindById(move.TypeId())

	typeIconFile, err = os.Open(fmt.Sprintf("./rawData/imgs/typeIcons/%s-icon.png", typeDto.DisplayName()))
	if err != nil {
		return err
	}
	typeIconImg, err = png.Decode(typeIconFile)
	if err != nil {
		return err
	}

	draw.Draw(img, image.Rect(45, 405, 85, 445), typeIconImg, image.Pt(0, 0), draw.Over)

	h := font.HintingNone
	nameDrawer := font.Drawer{
		Dst: img,
		Src: textColor,
		Face: truetype.NewFace(&ttf, &truetype.Options{
			Size:    24.0,
			DPI:     75.0,
			Hinting: h,
		}),
	}
	nameDrawer.Dot = fixed.Point26_6{
		X: fixed.I(90),
		Y: fixed.I(423),
	}
	nameDrawer.DrawString(move.Name())

	extraDrawer := font.Drawer{
		Dst: img,
		Src: textColor,
		Face: truetype.NewFace(&ttf, &truetype.Options{
			Size:    16.0,
			DPI:     75.0,
			Hinting: h,
		}),
	}
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(90),
		Y: fixed.I(441),
	}
	fastMove := pokemon.FastMove()
	fastAttacks := math.Ceil(-move.Energy() / fastMove.Energy())
	extraDrawer.DrawString(fmt.Sprintf("Fast Attacks: %.0f", fastAttacks))

	power := move.Power() * pokemon.GetStab(&move)
	powerLabel := fmt.Sprintf("Power: %.1f", power)
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(355+52) - extraDrawer.MeasureString(powerLabel)/2,
		Y: fixed.I(420),
	}
	extraDrawer.DrawString(powerLabel)

	energyLabel := fmt.Sprintf("Energy: %.0f", move.Energy())
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(355+52) - extraDrawer.MeasureString(energyLabel)/2,
		Y: fixed.I(441),
	}
	extraDrawer.DrawString(energyLabel)
	return nil
}

func applySecondaryCharge(img *image.RGBA, pokemon Pokemon, ttf truetype.Font, textColor *image.Uniform) error {
	var (
		typeIconFile *os.File
		typeIconImg  image.Image
		move         = pokemon.ChargeMoves()[1]
		err          error
	)
	_, typeDto := TYPES_DAO.FindById(move.TypeId())

	typeIconFile, err = os.Open(fmt.Sprintf("./rawData/imgs/typeIcons/%s-icon.png", typeDto.DisplayName()))
	if err != nil {
		return err
	}
	typeIconImg, err = png.Decode(typeIconFile)
	if err != nil {
		return err
	}

	draw.Draw(img, image.Rect(45, 455, 85, 495), typeIconImg, image.Pt(0, 0), draw.Over)

	h := font.HintingNone
	nameDrawer := font.Drawer{
		Dst: img,
		Src: textColor,
		Face: truetype.NewFace(&ttf, &truetype.Options{
			Size:    24.0,
			DPI:     75.0,
			Hinting: h,
		}),
	}
	nameDrawer.Dot = fixed.Point26_6{
		X: fixed.I(90),
		Y: fixed.I(473),
	}
	nameDrawer.DrawString(move.Name())

	extraDrawer := font.Drawer{
		Dst: img,
		Src: textColor,
		Face: truetype.NewFace(&ttf, &truetype.Options{
			Size:    16.0,
			DPI:     75.0,
			Hinting: h,
		}),
	}
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(90),
		Y: fixed.I(491),
	}
	fastMove := pokemon.FastMove()
	fastAttacks := math.Ceil(-move.Energy() / fastMove.Energy())
	extraDrawer.DrawString(fmt.Sprintf("Fast Attacks: %.0f", fastAttacks))

	power := move.Power() * pokemon.GetStab(&move)
	powerLabel := fmt.Sprintf("Power: %.1f", power)
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(355+52) - extraDrawer.MeasureString(powerLabel)/2,
		Y: fixed.I(470),
	}
	extraDrawer.DrawString(powerLabel)

	energyLabel := fmt.Sprintf("Energy: %.0f", move.Energy())
	extraDrawer.Dot = fixed.Point26_6{
		X: fixed.I(355+52) - extraDrawer.MeasureString(energyLabel)/2,
		Y: fixed.I(491),
	}
	extraDrawer.DrawString(energyLabel)
	return nil
}

func applyMatchups(img *image.RGBA, moveSetId int64, cup string) error {
	var (
		zeroLosses []int64
		zeroWins   []int64
		oneLosses  []int64
		oneWins    []int64
		twoLosses  []int64
		twoWins    []int64
		query      string
		rows       *sql.Rows
		err        error
	)

	query = `SELECT r.pokemon_id
FROM rankings r
INNER JOIN battle_simulations bs ON r.move_set_id = bs.ally_id AND bs.enemy_id = ?
WHERE r.cup = ?
  AND r.pokemon_rank IS NOT NULL
  AND bs.0v0 > 500
  AND bs.ally_id != bs.enemy_id
ORDER BY r.pokemon_rank ASC
LIMIT 5`
	rows, err = LIVE.Query(query, moveSetId, cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return err
	}
	for rows.Next() {
		var enemyId int64
		if rows.Scan(&enemyId) != nil {
			return err
		}
		zeroLosses = append(zeroLosses, enemyId)
	}
	err = rows.Close()
	if err != nil {
		return err
	}

	for i, id := range zeroLosses {
		pokemonFile, err := os.Open(fmt.Sprintf("./rawData/imgs/pokemon/%d.png", id))
		if err != nil {
			return err
		}
		pokemonImg, err := png.Decode(pokemonFile)
		if err != nil {
			return err
		}

		bounds := pokemonImg.Bounds()
		width := bounds.Size().X
		height := bounds.Size().Y
		widthMultiplier := 40.0 / float64(width)
		heightMultiplier := 40.0 / float64(height)
		multiplier := math.Min(widthMultiplier, heightMultiplier)
		newWidth := int(float64(width) * multiplier)
		newHeight := int(float64(height) * multiplier)
		x := 75 + i*30 + (40-newWidth)/2
		y := 535 + (40-newHeight)/2

		draw.BiLinear.Scale(img, image.Rect(x, y, x+newWidth, y+newHeight), pokemonImg, pokemonImg.Bounds(), draw.Over, nil)
	}

	query = `SELECT r.pokemon_id
FROM rankings r
INNER JOIN battle_simulations bs ON r.move_set_id = bs.ally_id AND bs.enemy_id = ?
WHERE r.cup = ?
  AND r.pokemon_rank IS NOT NULL
  AND bs.1v1 > 500
  AND bs.ally_id != bs.enemy_id
ORDER BY r.pokemon_rank ASC
LIMIT 5`
	rows, err = LIVE.Query(query, moveSetId, cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return err
	}
	for rows.Next() {
		var enemyId int64
		if rows.Scan(&enemyId) != nil {
			return err
		}
		oneLosses = append(oneLosses, enemyId)
	}
	err = rows.Close()
	if err != nil {
		return err
	}

	for i, id := range oneLosses {
		pokemonFile, err := os.Open(fmt.Sprintf("./rawData/imgs/pokemon/%d.png", id))
		if err != nil {
			return err
		}
		pokemonImg, err := png.Decode(pokemonFile)
		if err != nil {
			return err
		}

		bounds := pokemonImg.Bounds()
		width := bounds.Size().X
		height := bounds.Size().Y
		widthMultiplier := 40.0 / float64(width)
		heightMultiplier := 40.0 / float64(height)
		multiplier := math.Min(widthMultiplier, heightMultiplier)
		newWidth := int(float64(width) * multiplier)
		newHeight := int(float64(height) * multiplier)
		x := 75 + i*30 + (40-newWidth)/2
		y := 585 + (40-newHeight)/2

		draw.BiLinear.Scale(img, image.Rect(x, y, x+newWidth, y+newHeight), pokemonImg, pokemonImg.Bounds(), draw.Over, nil)
	}

	query = `SELECT r.pokemon_id
FROM rankings r
INNER JOIN battle_simulations bs ON r.move_set_id = bs.ally_id AND bs.enemy_id = ?
WHERE r.cup = ?
  AND r.pokemon_rank IS NOT NULL
  AND bs.2v2 > 500
  AND bs.ally_id != bs.enemy_id
ORDER BY r.pokemon_rank ASC
LIMIT 5`
	rows, err = LIVE.Query(query, moveSetId, cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return err
	}
	for rows.Next() {
		var enemyId int64
		if rows.Scan(&enemyId) != nil {
			return err
		}
		twoLosses = append(twoLosses, enemyId)
	}
	err = rows.Close()
	if err != nil {
		return err
	}

	for i, id := range twoLosses {
		pokemonFile, err := os.Open(fmt.Sprintf("./rawData/imgs/pokemon/%d.png", id))
		if err != nil {
			return err
		}
		pokemonImg, err := png.Decode(pokemonFile)
		if err != nil {
			return err
		}

		bounds := pokemonImg.Bounds()
		width := bounds.Size().X
		height := bounds.Size().Y
		widthMultiplier := 40.0 / float64(width)
		heightMultiplier := 40.0 / float64(height)
		multiplier := math.Min(widthMultiplier, heightMultiplier)
		newWidth := int(float64(width) * multiplier)
		newHeight := int(float64(height) * multiplier)
		x := 75 + i*30 + (40-newWidth)/2
		y := 635 + (40-newHeight)/2

		draw.BiLinear.Scale(img, image.Rect(x, y, x+newWidth, y+newHeight), pokemonImg, pokemonImg.Bounds(), draw.Over, nil)
	}

	query = `SELECT r.pokemon_id
FROM rankings r
INNER JOIN battle_simulations bs ON r.move_set_id = bs.ally_id AND bs.enemy_id = ?
WHERE r.cup = ?
  AND r.pokemon_rank IS NOT NULL
  AND bs.0v0 <
   500
  AND bs.ally_id != bs.enemy_id
ORDER BY r.pokemon_rank ASC
LIMIT 5`
	rows, err = LIVE.Query(query, moveSetId, cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return err
	}
	for rows.Next() {
		var enemyId int64
		if rows.Scan(&enemyId) != nil {
			return err
		}
		zeroWins = append(zeroWins, enemyId)
	}
	err = rows.Close()
	if err != nil {
		return err
	}

	for i, id := range zeroWins {
		pokemonFile, err := os.Open(fmt.Sprintf("./rawData/imgs/pokemon/%d.png", id))
		if err != nil {
			return err
		}
		pokemonImg, err := png.Decode(pokemonFile)
		if err != nil {
			return err
		}

		bounds := pokemonImg.Bounds()
		width := bounds.Size().X
		height := bounds.Size().Y
		widthMultiplier := 40.0 / float64(width)
		heightMultiplier := 40.0 / float64(height)
		multiplier := math.Min(widthMultiplier, heightMultiplier)
		newWidth := int(float64(width) * multiplier)
		newHeight := int(float64(height) * multiplier)
		x := 310 + i*30 + (40-newWidth)/2
		y := 535 + (40-newHeight)/2

		draw.BiLinear.Scale(img, image.Rect(x, y, x+newWidth, y+newHeight), pokemonImg, pokemonImg.Bounds(), draw.Over, nil)
	}

	query = `SELECT r.pokemon_id
FROM rankings r
INNER JOIN battle_simulations bs ON r.move_set_id = bs.ally_id AND bs.enemy_id = ?
WHERE r.cup = ?
  AND r.pokemon_rank IS NOT NULL
  AND bs.1v1 < 500
  AND bs.ally_id != bs.enemy_id
ORDER BY r.pokemon_rank ASC
LIMIT 5`
	rows, err = LIVE.Query(query, moveSetId, cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return err
	}
	for rows.Next() {
		var enemyId int64
		if rows.Scan(&enemyId) != nil {
			return err
		}
		oneWins = append(oneWins, enemyId)
	}
	err = rows.Close()
	if err != nil {
		return err
	}

	for i, id := range oneWins {
		pokemonFile, err := os.Open(fmt.Sprintf("./rawData/imgs/pokemon/%d.png", id))
		if err != nil {
			return err
		}
		pokemonImg, err := png.Decode(pokemonFile)
		if err != nil {
			return err
		}

		bounds := pokemonImg.Bounds()
		width := bounds.Size().X
		height := bounds.Size().Y
		widthMultiplier := 40.0 / float64(width)
		heightMultiplier := 40.0 / float64(height)
		multiplier := math.Min(widthMultiplier, heightMultiplier)
		newWidth := int(float64(width) * multiplier)
		newHeight := int(float64(height) * multiplier)
		x := 310 + i*30 + (40-newWidth)/2
		y := 585 + (40-newHeight)/2

		draw.BiLinear.Scale(img, image.Rect(x, y, x+newWidth, y+newHeight), pokemonImg, pokemonImg.Bounds(), draw.Over, nil)
	}

	query = `SELECT r.pokemon_id
FROM rankings r
INNER JOIN battle_simulations bs ON r.move_set_id = bs.ally_id AND bs.enemy_id = ?
WHERE r.cup = ?
  AND r.pokemon_rank IS NOT NULL
  AND bs.2v2 < 500
  AND bs.ally_id != bs.enemy_id
ORDER BY r.pokemon_rank ASC
LIMIT 5`
	rows, err = LIVE.Query(query, moveSetId, cup)
	defer func() {
		if rows != nil {
			_ = rows.Close()
		}
	}()
	if err != nil {
		return err
	}
	for rows.Next() {
		var enemyId int64
		if rows.Scan(&enemyId) != nil {
			return err
		}
		twoWins = append(twoWins, enemyId)
	}
	err = rows.Close()
	if err != nil {
		return err
	}

	for i, id := range twoWins {
		pokemonFile, err := os.Open(fmt.Sprintf("./rawData/imgs/pokemon/%d.png", id))
		if err != nil {
			return err
		}
		pokemonImg, err := png.Decode(pokemonFile)
		if err != nil {
			return err
		}

		bounds := pokemonImg.Bounds()
		width := bounds.Size().X
		height := bounds.Size().Y
		widthMultiplier := 40.0 / float64(width)
		heightMultiplier := 40.0 / float64(height)
		multiplier := math.Min(widthMultiplier, heightMultiplier)
		newWidth := int(float64(width) * multiplier)
		newHeight := int(float64(height) * multiplier)
		x := 310 + i*30 + (40-newWidth)/2
		y := 635 + (40-newHeight)/2

		draw.BiLinear.Scale(img, image.Rect(x, y, x+newWidth, y+newHeight), pokemonImg, pokemonImg.Bounds(), draw.Over, nil)
	}

	return nil
}
