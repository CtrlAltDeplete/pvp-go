package db

type TypeMultipliersDao struct{}

func (dao *TypeMultipliersDao) FindByIds(receivingType, actingType int64) (error, float64) {
	var (
		e          error
		multiplier float64
		query      = "SELECT multiplier " +
			"FROM pvpgo.type_multipliers " +
			"WHERE acting_type = ? " +
			"AND receiving_type = ?"
	)
	e = LIVE.QueryRow(query, actingType, receivingType).Scan(&multiplier)
	return e, multiplier
}

func (dao *TypeMultipliersDao) Create(receivingType, actingType int64, multiplier float64) {
	var (
		e     error
		query = "INSERT INTO pvpgo.type_multipliers (receiving_type, acting_type, multiplier) " +
			"VALUES (?, ?, ?)"
	)
	_, e = LIVE.Exec(query, receivingType, actingType, multiplier)
	CheckError(e)
}

func (dao *TypeMultipliersDao) FindOrCreate(receivingType, actingType int64, multiplier float64) {
	var (
		err error
	)
	err, _ = dao.FindByIds(receivingType, actingType)
	if err != nil {
		dao.Create(receivingType, actingType, multiplier)
	}
}
