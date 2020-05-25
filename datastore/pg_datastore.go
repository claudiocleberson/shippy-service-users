package datastore

import (
	"context"

	"github.com/claudiocleberson/shippy-service-users/models"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type DatastoreClient interface {
	Create(context.Context, *models.User) error
	Get(context.Context, string) (*models.User, error)
	GetAll(context.Context) (models.Users, error)
	GetByEmailAndPassword(context.Context, *models.User) (*models.User, error)
	Auth(context.Context, *models.User) error
	ValidateToken(context.Context, *models.Token) (bool, error)
}

type datastoreClient struct {
	db *gorm.DB
}

func NewDatastoreClient(dbstring string) DatastoreClient {

	db, err := gorm.Open("postgres", dbstring)
	if err != nil {
		panic(err)
	}

	db.AutoMigrate(&models.User{}, &models.Token{})

	return &datastoreClient{
		db: db,
	}
}

func (d *datastoreClient) Create(ctx context.Context, user *models.User) error {

	if err := d.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func (d *datastoreClient) Get(ctx context.Context, id string) (*models.User, error) {
	var user *models.User
	user.UserID = id
	if err := d.db.First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (d *datastoreClient) GetByEmailAndPassword(ctx context.Context, user *models.User) (*models.User, error) {

	if err := d.db.Where("email =? AND password = ?", user.Email, user.Password).First(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (d *datastoreClient) GetAll(ctx context.Context) (models.Users, error) {

	var users models.Users
	if err := d.db.Find(&users).Error; err != nil {
		return nil, err
	}
	return users, nil
}

func (d *datastoreClient) Auth(context.Context, *models.User) error {
	return nil
}

func (d *datastoreClient) ValidateToken(context.Context, *models.Token) (bool, error) {
	return false, nil
}
