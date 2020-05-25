package handlers

import (
	"context"

	"github.com/claudiocleberson/shippy-service-users/models"
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/claudiocleberson/shippy-service-users/repository"
)

type UserServiceHandler interface {
	Create(context.Context, *pb.User, *pb.Response) error
	Get(context.Context, *pb.User, *pb.Response) error
	GetAll(context.Context, *pb.Request, *pb.Response) error
	Auth(context.Context, *pb.User, *pb.Token) error
	ValidateToken(context.Context, *pb.Token, *pb.Token) error
}

func NewUserserviceHandler(userRepo repository.UserRepository) UserServiceHandler {
	return userServiceHandler{
		userRepo: userRepo,
	}
}

type userServiceHandler struct {
	userRepo repository.UserRepository
	tokeRepo repository.AuthRepository
}

func (s userServiceHandler) Create(ctx context.Context, req *pb.User, res *pb.Response) error {

	if err := s.userRepo.Create(ctx, models.MarshalUser(req)); err != nil {
		return err
	}

	//Todo - return the user with id created on DB
	res.User = req

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

	_, err := s.userRepo.GetByEmailAndPassword(ctx, models.MarshalUser(req))
	if err != nil {
		return err
	}

	res.Token = "testing_token"

	return nil
}

func (s userServiceHandler) ValidateToken(ctx context.Context, req *pb.Token, res *pb.Token) error {

	return nil
}
