package repository

import (
	"context"

	"github.com/claudiocleberson/shippy-service-users/models"
)

type AuthRepository interface {
	Auth(context.Context, *models.User) (*models.Token, error)
	ValidateToken(context.Context, *models.Token) (bool, error)
}

type authRepository struct{}

func (repo *authRepository) Auth(ctx context.Context, user *models.User) (*models.Token, error) {
	return nil, nil
}

func (repo *authRepository) ValidateToken(ctx context.Context, token *models.Token) (bool, error) {

	return false, nil
}
