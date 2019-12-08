package dtos

type ApiRankingDto struct {
	Name         string          `json:"name"`
	Id           int64           `json:"id"`
	RelativeRank int64           `json:"relativeRank"`
	MoveSets     []ApiMoveSetDto `json:"moveSets"`
}

type ApiMoveSetDto struct {
	Id                  int64       `json:"id"`
	AbsoluteRank        float64     `json:"absoluteRank"`
	FastMove            ApiMoveDto  `json:"fastMove"`
	PrimaryChargeMove   ApiMoveDto  `json:"primaryChargeMove"`
	SecondaryChargeMove *ApiMoveDto `json:"secondaryChargeMove"`
}

type ApiMoveDto struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
