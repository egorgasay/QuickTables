package repository

func (s Storage) AddDB(strCon string, owner string, driver string) error {
	query := "INSERT INTO userDBs (connStr, owner, driver) VALUES (?, ?, ?)"
	_, err := s.DB.Exec(query, strCon, owner, driver)
	return err
}
