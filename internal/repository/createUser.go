package repository

func (s Storage) CreateUser(username string, password string) error {
	query := "INSERT INTO Users (Name, Role, password, Nickname) VALUES (?, 'User', ?, ?)"
	_, err := s.DB.Exec(query, username, password, "No Name")
	return err
}
