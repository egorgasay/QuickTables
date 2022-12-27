package repository

import "log"

func (s Storage) GetAllDBs(username string) [][]string {
	rows, err := s.DB.Query("SELECT dbName,driver FROM userDBs WHERE owner = ?", username)
	if err != nil {
		log.Println(err)
	}

	names := make([][]string, 0, 5)

	for rows.Next() {
		name, driver := "", ""
		rows.Scan(&name, &driver)
		names = append(names, []string{name, driver})
	}

	return names
}
