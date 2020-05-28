package handlers

import (
	"context"
	"errors"

	"github.com/claudiocleberson/shippy-service-users/models"
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/claudiocleberson/shippy-service-users/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceHandler interface {
	Create(context.Context, *pb.User, *pb.Response) error
	Get(context.Context, *pb.User, *pb.Response) error
	GetAll(context.Context, *pb.Request, *pb.Response) error
	Auth(context.Context, *pb.User, *pb.Token) error
	ValidateToken(context.Context, *pb.Token, *pb.Token) error
}

func NewUserserviceHandler(userRepo repository.UserRepository, tokenRepo repository.AuthRepository) UserServiceHandler {
	return userServiceHandler{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
	}
}

type userServiceHandler struct {
	userRepo  repository.UserRepository
	tokenRepo repository.AuthRepository
}

func (s userServiceHandler) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	hashedPass, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	req.Password = string(hashedPass)
	if result, err := s.userRepo.Create(ctx, models.MarshalUser(req)); err != nil {
		return err
	} else {
		res.User = models.UnmarshalUser(result)
	}

	//Todo - return the user with id created on DB

	return nil
}

func (s userServiceHandler) Get(ctx context.Context, req *pb.User, res *pb.Response) error {
	user, err := s.userRepo.Get(ctx, req.Id)
	if err != nil {
		return err
	}

	res.User = models.UnmarshalUser(user)

	return nil
}

func (s userServiceHandler) GetAll(ctx context.Context, req *pb.Request, res *pb.Response) error {

	users, err := s.userRepo.GetAll(ctx)
	if err != nil {
		return err
	}

	res.Users = models.UnmarshalUserCollection(users)

	return nil
}

func (s userServiceHandler) Auth(ctx context.Context, req *pb.User, res *pb.Token) error {

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

func (s userServiceHandler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	return nil
}
