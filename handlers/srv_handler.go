package handlers

import (
	"context"
	"errors"
	"fmt"

	"github.com/claudiocleberson/shippy-service-users/models"
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/claudiocleberson/shippy-service-users/repository"
	"github.com/micro/go-micro"

	"golang.org/x/crypto/bcrypt"
)

const (
	topic = "user.created"
)

type UserServiceHandler interface {
	Create(context.Context, *pb.User, *pb.Response) error
	Get(context.Context, *pb.User, *pb.Response) error
	GetAll(context.Context, *pb.Request, *pb.Response) error
	Auth(context.Context, *pb.User, *pb.Token) error
	ValidateToken(context.Context, *pb.Token, *pb.Token) error
}

func NewUserserviceHandler(userRepo repository.UserRepository,
	tokenRepo repository.AuthRepository,
	publisher micro.Publisher) UserServiceHandler {
	return &userServiceHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		publisher: publisher,
	}
}

type userServiceHandler struct {
	userRepo  repository.UserRepository
	tokenRepo repository.AuthRepository
	publisher micro.Publisher
}

func (s *userServiceHandler) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New(fmt.Sprintf("error hashing password: %v", err))
	}

	req.Password = string(hashedPass)
	if result, err := s.userRepo.Create(ctx, models.MarshalUser(req)); err != nil {
		return errors.New(fmt.Sprintf("error creating user: %v", err))
	} else {
		res.User = models.UnmarshalUser(result)
	}

	//Publish the user created event
	if err := s.publisher.Publish(ctx, res.User); err != nil {
		return errors.New(fmt.Sprintf("error publishing event: %v", err))
	}

	return nil
}

func (s *userServiceHandler) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.userRepo.Get(ctx, req.Id)
	if err != nil {
		return err
	}

	res.User = models.UnmarshalUser(user)

	return nil
}

func (s *userServiceHandler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {

	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	res.Users = models.UnmarshalUserCollection(users)

	return nil
}

func (s *userServiceHandler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {

	mUser := models.MarshalUser(req)
	user, err := s.userRepo.GetByEmail(ctx, mUser)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return errors.New("Password does not match!")
	}

	token, err := s.tokenRepo.Auth(ctx, user)
	if err != nil {
		return err
	}

	res.Token = token.Token

	return nil
}

func (s *userServiceHandler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	//Decode the token
	_, err := s.tokenRepo.ValidateToken(ctx, models.MarshalToken(req))
	if err != nil {
		return err
	}

	res.Token = req.Token
	res.Valid = true
	res.Errors = nil

	return nil
}

// func (s *userServiceHandler) publishEvent(user *pb.User) error {

// 	//Marshal to JSON string
// 	body, err := json.Marshal(user)
// 	if err != nil {
// 		return err
// 	}

// 	// Crete a broker message
// 	msg := &broker.Message{
// 		Header: map[string]string{
// 			"id": user.Id,
// 		},
// 		Body: body,
// 	}

// 	//Publish message to broker
// 	if err := s.PubSub.Publish(topic, msg); err != nil {
// 		log.Printf("[PUB] failed: %v", err)
// 		return err
// 	}
// 	return nil
// }
