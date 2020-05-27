package main

import (
	"log"
	"os"
	"strings"

	"github.com/claudiocleberson/shippy-service-users/datastore"
	"github.com/claudiocleberson/shippy-service-users/handlers"
	pb "github.com/claudiocleberson/shippy-service-users/proto/users"
	"github.com/claudiocleberson/shippy-service-users/repository"
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
