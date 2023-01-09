package repository

func (s Storage) AddDB(dbName, strCon, owner, driver, docker string) error {
	query := "INSERT INTO userDBs (connStr, owner, driver, dbName, docker) VALUES (?, ?, ?, ?, ?)"
	_, err := s.DB.Exec(query, strCon, owner, driver, dbName, docker)
	return err
}
