package rawData

import (
	"math"
	"testing"
)

func TestFillTypeChartDto(t *testing.T) {
	var (
		typeChartDto = typeChartDto{
			Name: `Grass\/Poison`,
			FieldTypeAdvantage: `
###
###
###
	<tr>
		<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/electric.gif"/></td>
		<td>Takes <span class="type-resist-value type-resist-value-62.5">62.5%</span> damage</td>
	</tr>
###
	<tr>
		<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/fairy.gif"/></td>
		<td>Takes <span class="type-resist-value type-resist-value-62.5">62.5%</span> damage</td>
	</tr>
###
	<tr>
		<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/fighting.gif"/></td>
		<td>Takes <span class="type-resist-value type-resist-value-62.5">62.5%</span> damage</td>
	</tr>
###
<tr>
	<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/fire.gif"/></td>
	<td>Takes <span class="type-weak-value type-weak-value-160">160%</span> damage</td>
</tr>
###
<tr>
	<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/flying.gif"/></td>
	<td>Takes <span class="type-weak-value type-weak-value-160">160%</span> damage</td>
</tr>
###
###
	<tr>
		<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/grass.gif"/></td>
		<td>Takes <span class="type-resist-value type-resist-value-39.1">39.1%</span> damage</td>
	</tr>
###
###
<tr>
	<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/ice.gif"/></td>
	<td>Takes <span class="type-weak-value type-weak-value-160">160%</span> damage</td>
</tr>
###
###
###
<tr>
	<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/psychic.gif"/></td>
	<td>Takes <span class="type-weak-value type-weak-value-160">160%</span> damage</td>
</tr>
###
###
###
	<tr>
		<td class="type-img-cell"><img src="/sites/pokemongo/files/2016-07/water.gif"/></td>
		<td>Takes <span class="type-resist-value type-resist-value-62.5">62.5%</span> damage</td>
	</tr>
`,
			types:       nil,
			multipliers: nil,
		}
		expectedTypes       = []string{"grass", "poison"}
		expectedMultipliers = map[string]float64{
			"normal":   1,
			"fire":     1.60,
			"fighting": 0.625,
			"water":    0.625,
			"flying":   1.60,
			"grass":    0.391,
			"poison":   1,
			"electric": 0.625,
			"ground":   1,
			"psychic":  1.60,
			"rock":     1,
			"ice":      1.60,
			"bug":      1,
			"dragon":   1,
			"ghost":    1,
			"dark":     1,
			"steel":    1,
			"fairy":    0.625,
		}
	)

	FillTypeChartDto(&typeChartDto)
	compareTypeChartDto(typeChartDto, expectedTypes, expectedMultipliers, "FillTypeChartDto", t)
}

func compareTypeChartDto(dto typeChartDto, expectedTypes []string, expectedMultipliers map[string]float64, functionName string, t *testing.T) {
	if len(expectedTypes) != len(dto.types) {
		t.Errorf("%s() expected %d types; got %d", functionName, len(expectedTypes), len(dto.types))
	}

	for i := 0; i < len(expectedTypes); i++ {
		if expectedTypes[i] != dto.types[i] {
			t.Errorf("%s() expected %s as %d type; got %s", functionName, expectedTypes[i], i, dto.types[i])
		}
	}

	for key, value := range expectedMultipliers {
		if math.Abs(value-dto.multipliers[key]) > 0.01 {
			t.Errorf("%s() expected %s multiplier of %f; got %f", functionName, key, value, dto.multipliers[key])
		}
	}
}
