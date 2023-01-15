package repository

func (s Storage) BindPort(port string) error {
	prep, err := s.DB.Prepare("INSERT INTO Ports VALUES (?)")
	if err != nil {
		return err
	}

	_, err = prep.Exec(port)
	return err
}
