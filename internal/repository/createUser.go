package repository

func (s Storage) CreateUser(username string, password string) error {
	query := "INSERT INTO Users (Name, Role, password) VALUES (?, 'user', ?)"
	_, err := s.DB.Exec(query, username, password)
	return err
}
