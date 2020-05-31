package main

import (
	"log"
	"os"
	"strings"

	"github.com/claudiocleberson/shippy-service-users/datastore"
	"github.com/claudiocleberson/shippy-service-users/handlers"
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/claudiocleberson/shippy-service-users/repository"
	"github.com/claudiocleberson/shippy-service-users/services"
	"github.com/micro/go-micro"
)

const (
	dbHost    = "DB_HOST"
	dbDefault = "host=127.0.0.1 port=5432 user=example dbname=users password=example sslmode=disable"
)

func main() {

	dbString := os.Getenv(dbHost)

	if strings.TrimSpace(dbString) == "" {
		log.Println("missing or empty DB_HOST environment variable...")
		log.Printf("connecting on dev database : %s", dbDefault)
		dbString = dbDefault
	}

	dbClient := datastore.NewDatastoreClient(dbString)
	tokenRepo := repository.NewAuthRepository(services.NewTokenService())
	repo := repository.NewUserRepository(dbClient)

	srv := micro.NewService(
		micro.Name("shippy.service.users"),
	)
	srv.Init()

	// Get instance of the broker using our defaults
	pubsub := srv.Server().Options().Broker

	userServiceHandler := handlers.NewUserserviceHandler(repo, tokenRepo, pubsub)

	pb.RegisterUserServiceHandler(srv.Server(), userServiceHandler)

	if err := srv.Run(); err != nil {
		panic(err)
	}

}
