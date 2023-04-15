package repository

func (s Storage) GetAllDBs(username string) ([][]string, error) {
	rows, err := s.DB.Query("SELECT dbName,driver FROM userDBs WHERE owner = ?", username)
	if err != nil {
		return nil, err
	}

	names := make([][]string, 0, 5)

	for rows.Next() {
		name, driver := "", ""
		err = rows.Scan(&name, &driver)
		if err != nil {
			return nil, err
		}
		names = append(names, []string{name, driver})
	}

	return names, nil
}
