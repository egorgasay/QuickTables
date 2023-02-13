package usecase

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

func (uc *UseCase) GetAllDBs(owner string) [][]string {
	return uc.service.DB.GetAllDBs(owner)
}
