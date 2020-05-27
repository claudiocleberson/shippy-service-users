package datastore

import (
	"context"
	"log"
	"time"

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

var (
	retry int
)

type datastoreClient struct {
	db *gorm.DB
}

func NewDatastoreClient(dbstring string) DatastoreClient {

	db := connectDatabaseCluster(dbstring)

	db.AutoMigrate(&models.User{}, &models.Token{})

	return &datastoreClient{
		db: db,
	}
}

func connectDatabaseCluster(dbstring string) *gorm.DB {

	log.Println("Connecting database....")

	db, err := gorm.Open("postgres", dbstring)
	if err != nil {
		if retry >= 3 {
			panic(err)
		}
		retry = retry + 1
		time.Sleep(time.Second * 2)
		connectDatabaseCluster(dbstring)
	}

	log.Println("Database conneted....")

	return db
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
