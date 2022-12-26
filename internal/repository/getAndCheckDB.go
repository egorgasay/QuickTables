package repository

func (s Storage) CheckDB(owner string) bool {
	row := s.DB.QueryRow("SELECT count(*) FROM userDBs WHERE owner = ?", owner)
	var count int
	row.Scan(&count)
	if count > 0 {
		return true
	}
	return false
}

func (s Storage) GetDB(owner string) (connStr, driver string) {
	row := s.DB.QueryRow("SELECT connstr, driver FROM userDBs WHERE owner = ?", owner)
	row.Scan(&connStr, &driver)
	return connStr, driver
}
