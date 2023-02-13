package usecase

import "quicktables/internal/repository"

func (uc *UseCase) GetProfile(username string) (*repository.UserStats, error) {
	queries, err := uc.service.DB.GetUserStats(username)
	if err != nil {
		return nil, err
	}

	return queries, nil
}
