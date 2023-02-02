package usecase

import "quicktables/internal/repository"

func (uc UseCase) GetHistory(username string) ([]repository.QueryInfo, error) {
	dbName, err := uc.userDBs.GetDBName(username)
	if err != nil {
		return nil, err
	}

	queries, err := uc.Service.DB.GetQueries(username, dbName)
	if err != nil {
		return nil, err
	}

	return queries, nil
}
