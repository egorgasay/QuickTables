package repository

func (s Storage) ChangeNick(username, nick string) error {
	query := `
	UPDATE Users 
	SET Nickname = ?
	WHERE Name = ?;
`
	stmt, err := s.DB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(nick, username)
	if err != nil {
		return err
	}

	return nil
}
