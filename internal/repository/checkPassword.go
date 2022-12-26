package repository

func (s Storage) CheckPassword(username string, password string) bool {
	query := "SELECT count(*) FROM Users WHERE Name = ? AND Password = ?"
	var res int

	row := s.DB.QueryRow(query, username, password)
	err := row.Scan(&res)

	if err != nil || res < 1 {
		return false
	}

	return true
}
