package orm

func CalcPageLimit(page, limit int) (int, int) {
	return limit, limit * (page - 1)
}
