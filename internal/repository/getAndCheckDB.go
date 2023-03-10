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

func (s Storage) GetDB(owner string) (connStr, driver, dbName string) {
	row := s.DB.QueryRow("SELECT connStr, driver, dbName FROM userDBs WHERE owner = ? LIMIT 1", owner)
	row.Scan(&connStr, &driver, &dbName)
	return dbName, connStr, driver
}

func (s Storage) GetDBInfobyName(owner, name string) (connStr, driver, docker string) {
	row := s.DB.QueryRow("SELECT connStr, driver, docker FROM userDBs WHERE owner = ? AND  dbName = ? LIMIT 1", owner, name)
	row.Scan(&connStr, &driver, &docker)
	return connStr, driver, docker
}
