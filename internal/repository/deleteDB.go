package repository

func (s Storage) DeleteDB(username, dbName string) error {
	_, err := s.DB.Exec("DELETE FROM userDBs WHERE owner = ? and dbName = ?",
		username, dbName)
	return err
}
