package usecase

func (uc *UseCase) BindPort(port string) error {
	return uc.service.DB.BindPort(port)
}

func (uc *UseCase) ChangeNick(username string, nick string) error {
	return uc.service.DB.ChangeNick(username, nick)
}

func (uc *UseCase) ChangePassword(username string, oldPassword string, newPassword string) error {
	return uc.service.DB.ChangePassword(username, oldPassword, newPassword)
}

func (uc *UseCase) CreateUser(username string, password string) error {
	return uc.service.DB.CreateUser(username, password)
}
