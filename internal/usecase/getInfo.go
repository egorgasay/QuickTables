package usecase

import "quicktables/internal/repository"

func (uc *UseCase) GetVendorAndName(username string) (vendor, currentDB string, err error) {
	vendor, err = uc.userDBs.GetDBVendor(username)
	if err != nil {
		return "", "", err
	}

	currentDB, err = uc.userDBs.GetDBName(username)
	if err != nil {
		return "", "", err
	}

	return vendor, currentDB, nil
}

func (uc *UseCase) GetAllDBs(owner string) ([][]string, error) {
	return uc.service.DB.GetAllDBs(owner)
}

func (uc *UseCase) GetHistory(username string) ([]repository.QueryInfo, error) {
	dbName, err := uc.userDBs.GetDBName(username)
	if err != nil {
		return nil, err
	}

	queries, err := uc.service.DB.GetQueries(username, dbName)
	if err != nil {
		return nil, err
	}

	return queries, nil
}

func (uc *UseCase) GetProfile(username string) (*repository.UserStats, error) {
	queries, err := uc.service.DB.GetUserStats(username)
	if err != nil {
		return nil, err
	}

	return queries, nil
}
