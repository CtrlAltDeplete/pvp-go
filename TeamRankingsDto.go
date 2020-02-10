package main

type TeamRankingDto struct {
	id          int64
	cup         string
	allyOneId   int64
	allyTwoId   int64
	allyThreeId int64
	score       float64
}

func (t *TeamRankingDto) Id() int64 {
	return t.id
}

func (t *TeamRankingDto) SetId(id int64) {
	t.id = id
}

func (t *TeamRankingDto) Cup() string {
	return t.cup
}

func (t *TeamRankingDto) SetCup(cup string) {
	t.cup = cup
}

func (t *TeamRankingDto) AllyOneId() int64 {
	return t.allyOneId
}

func (t *TeamRankingDto) SetAllyOneId(allyOneId int64) {
	t.allyOneId = allyOneId
}

func (t *TeamRankingDto) AllyTwoId() int64 {
	return t.allyTwoId
}

func (t *TeamRankingDto) SetAllyTwoId(allyTwoId int64) {
	t.allyTwoId = allyTwoId
}

func (t *TeamRankingDto) AllyThreeId() int64 {
	return t.allyThreeId
}

func (t *TeamRankingDto) SetAllyThreeId(allyThreeId int64) {
	t.allyThreeId = allyThreeId
}

func (t *TeamRankingDto) Score() float64 {
	return t.score
}

func (t *TeamRankingDto) SetScore(score float64) {
	t.score = score
}
