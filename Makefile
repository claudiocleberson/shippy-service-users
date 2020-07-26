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
	# docker run 	-p 50053:50051 \
	# 		--network=host \ 
	# 		-e MICRO_ADDRESS=:50051 \  
	# 		-e MICRO_REGISTRY=mdns \
	# 		$(PROJECTNAME) 

	docker run -p  50053:50051 \
				--network=host \
				-e DB_HOST="host=datastore2 port=5432 user=example dbname=users password=example sslmode=disable" \
	 		    -e MICRO_ADDRESS=:50051\
			    -e MICRO_REGISTRY=etcd\
				-e MICRO_REGISTRY_ADDRESS=etcd-server:2379 \
			    $(PROJECTNAME)