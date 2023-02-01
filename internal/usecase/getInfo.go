package usecase

func (uc UseCase) GetVendorAndName(username string) (string, string, error) {
	vendor, err := uc.userDBs.GetDBVendor(username)
	if err != nil {
		return "", "", err
	}

	currentDB, err := uc.userDBs.GetDBName(username)
	if err != nil {
		return "", "", err
	}

	return vendor, currentDB, nil
}
