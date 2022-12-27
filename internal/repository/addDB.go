package repository

func (s Storage) AddDB(dbName, strCon, owner, driver string) error {
	query := "INSERT INTO userDBs (connStr, owner, driver, dbName) VALUES (?, ?, ?, ?)"
	_, err := s.DB.Exec(query, strCon, owner, driver, dbName)
	return err
}
