package main

import (
	"os"
	"strings"

	"github.com/claudiocleberson/shippy-service-users/datastore"
	"github.com/claudiocleberson/shippy-service-users/handlers"
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/claudiocleberson/shippy-service-users/repository"
	"github.com/micro/go-micro"
)

const (
	port      = ":50053"
	dbHost    = "DB_HOST"
	dbDefault = "host=127.0.0.1 port=5432 user=example dbname=users password=example sslmode=disable"
)

func main() {

	dbString := os.Getenv(dbDefault)
	if strings.TrimSpace(dbString) == "" {
		//panic("missing or empty DB_HOST environment variable...")
	}

	dbClient := datastore.NewDatastoreClient(dbDefault)
	repo := repository.NewUserRepository(dbClient)

	srv := micro.NewService(
		micro.Name("shippy-service-users"),
	)
	srv.Init()

	userServiceHandler := handlers.NewUserserviceHandler(repo)

	pb.RegisterUserServiceHandler(srv.Server(), userServiceHandler)

	if err := srv.Run(); err != nil {
		panic(err)
	}

}
