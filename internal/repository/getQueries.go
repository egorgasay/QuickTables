package repository

import (
	"strconv"
	"strings"
	"time"
)

type QueryInfo struct {
	Query    string
	Date     string
	ExecTime string
	Status   int
	Height   int
}

func (s Storage) GetQueries(username, dbName string) ([]QueryInfo, error) {
	qiArr := make([]QueryInfo, 0, 10)
	query := `SELECT Query, Date, ExecutionTime, Status FROM historyOfQueries
				WHERE Author = ? AND DBName = ?`

	rows, err := s.DB.Query(query, username, dbName)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		qi := QueryInfo{}
		rows.Scan(&qi.Query, &qi.Date, &qi.ExecTime, &qi.Status)
		date, _ := strconv.Atoi(qi.Date)

		qi.Date = time.Unix(int64(date), 0).String()
		qi.Height = (len(strings.Split(qi.Query, "\n")) + 3) * 40
		qiArr = append(qiArr, qi)
	}

	for i, j := 0, len(qiArr)-1; i < j; i, j = i+1, j-1 {
		qiArr[i], qiArr[j] = qiArr[j], qiArr[i]
	}

	return qiArr, nil
}
