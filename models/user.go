package models

import (
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
)

type User struct {
	gorm.Model
	UserID   string `json:"user_id"`
	Name     string `json:"name"`
	Company  string `json:"company"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

//Create a UUID for the user ID.
func (u *User) BeforeCreate(scope *gorm.Scope) error {
	uuid := uuid.NewUUID()
	return scope.SetColumn("user_id", uuid.String())
}

type Users []*User

func MarshalUser(user *pb.User) *User {
	return &User{
		UserID:   user.Id,
		Name:     user.Name,
		Company:  user.Company,
		Email:    user.Email,
		Password: user.Password,
	}
}

func UnmarshalUser(user *User) *pb.User {
	return &pb.User{
		Id:       user.UserID,
		Name:     user.Name,
		Company:  user.Company,
		Email:    user.Email,
		Password: user.Password,
	}
}

func MarshalUserCollection(users []*pb.User) Users {
	var collection Users
	for _, u := range users {
		collection = append(collection, MarshalUser(u))
	}
	return collection
}

func UnmarshalUserCollection(users Users) []*pb.User {

	collection := make([]*pb.User, len(users))
	for index, u := range users {
		collection[index] = UnmarshalUser(u)
	}
	return collection
}
