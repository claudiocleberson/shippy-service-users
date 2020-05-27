PROJECTNAME=$(shell basename "$(PWD)")


proto-file:
	#Build proto file
	protoc -I. --go_out=plugins=micro:. proto/users/user.proto

build:
	#Build locally our app
	env CGO_ENABLED=0  GOOS=linux go build -a -installsuffix cgo -o ./builds/$(PROJECTNAME)

	#Building container
	docker build -t $(PROJECTNAME) .
run:
	#Running docker container
	docker run --network=host -p 50053:50051 -e MICRO_ADDRESS=:50051 -e MICRO_SERVER_ADDRESS=127.0.0.1:8080 -e MICRO_REGISTRY=mdns $(PROJECTNAME) 