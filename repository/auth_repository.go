package repository

import (
	"context"

	"github.com/claudiocleberson/shippy-service-users/models"
	"github.com/claudiocleberson/shippy-service-users/services"
)

type AuthRepository interface {
	Auth(context.Context, *models.User) (*models.Token, error)
	ValidateToken(context.Context, *models.Token) (bool, error)
}

type authRepository struct {
	tokenService services.TokenService
}

func NewAuthRepository(tService services.TokenService) AuthRepository {
	return &authRepository{
		tokenService: tService,
	}
}

func (repo *authRepository) Auth(ctx context.Context, user *models.User) (*models.Token, error) {
	token, err := repo.tokenService.Encode(user)
	if err != nil {
		return nil, err
	}

	t := &models.Token{
		Token: token,
		Valid: true,
	}

	return t, nil
}

func (repo *authRepository) ValidateToken(ctx context.Context, token *models.Token) (bool, error) {

	return false, nil
}
