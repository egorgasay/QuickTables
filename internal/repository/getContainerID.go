package repository

func (s Storage) GetContainerID(username, dbName string) (string, error) {
	prepare, err := s.DB.Prepare("SELECT docker FROM userDBs WHERE owner = ? AND dbName = ?")
	if err != nil {
		return "", err
	}

	var id string
	err = prepare.QueryRow(username, dbName).Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}
