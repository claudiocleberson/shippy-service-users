package repository

import (
	"context"

	"github.com/claudiocleberson/shippy-service-users/datastore"
	"github.com/claudiocleberson/shippy-service-users/models"
)

type UserRepository interface {
	Create(context.Context, *models.User) error
	Get(context.Context, string) (*models.User, error)
	GetAll(context.Context) (models.Users, error)
	GetByEmailAndPassword(context.Context, *models.User) (*models.User, error)
}

func NewUserRepository(client datastore.DatastoreClient) UserRepository {
	return &userRepository{
		dbClient: client,
	}
}

type userRepository struct {
	dbClient datastore.DatastoreClient
}

func (repo *userRepository) Create(ctx context.Context, user *models.User) error {

	err := repo.dbClient.Create(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

func (repo *userRepository) Get(ctx context.Context, id string) (*models.User, error) {
	user, err := repo.dbClient.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (repo *userRepository) GetAll(ctx context.Context) (models.Users, error) {
	users, err := repo.dbClient.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (repo *userRepository) GetByEmailAndPassword(ctx context.Context, user *models.User) (*models.User, error) {
	user, err := repo.dbClient.GetByEmailAndPassword(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}
