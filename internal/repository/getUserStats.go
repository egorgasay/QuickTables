package repository

import (
	"fmt"
)

type UserStats struct {
	Nickname     string
	Role         string
	QueriesCount int
	Successful   int
	Percent      string
}

func (s Storage) GetUserStats(username string) (*UserStats, error) {
	queryCount := `SELECT DISTINCT count(*) FROM historyOfQueries WHERE Author=?`
	queryCountSuccessful := `
	SELECT DISTINCT count(*) FROM historyOfQueries 
	                WHERE Author=?
	                AND status=1`
	queryGetRoleAndNick := `SELECT Role, Nickname FROM Users WHERE Name = ?`

	prepareCount, err := s.DB.Prepare(queryCount)
	if err != nil {
		return nil, err
	}

	prepareSuccessful, err := s.DB.Prepare(queryCountSuccessful)
	if err != nil {
		return nil, err
	}

	prepareRole, err := s.DB.Prepare(queryGetRoleAndNick)
	if err != nil {
		return nil, err
	}

	var us = &UserStats{}

	prepareSuccessful.QueryRow(username).Scan(&us.Successful)
	prepareCount.QueryRow(username).Scan(&us.QueriesCount)
	prepareRole.QueryRow(username).Scan(&us.Role, &us.Nickname)

	us.Percent = fmt.Sprint(uint8(float64(us.Successful) / float64(us.QueriesCount) * 100.0))

	return us, nil
}
