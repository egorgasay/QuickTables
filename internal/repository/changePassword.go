package repository

import "errors"

func (s Storage) ChangePassword(username, oldPassword, newPassword string) error {
	prepareCurrentPassword, err := s.DB.Prepare(
		`SELECT password FROM Users WHERE Name = ?`)
	if err != nil {
		return err
	}

	var currentPassword string
	prepareCurrentPassword.QueryRow(username).Scan(&currentPassword)
	if currentPassword != oldPassword {
		return errors.New("Wrong current password")
	}

	prepareUpdatePassword, err := s.DB.Prepare(
		`UPDATE Users SET Password = ? WHERE Name = ?`)
	if err != nil {
		return err
	}

	_, err = prepareUpdatePassword.Exec(newPassword, username)
	if err != nil {
		return err
	}

	return nil
}
