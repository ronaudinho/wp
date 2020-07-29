### quick start
```
git clone https://github.com/ronaudinho/wp
cd wp/cmd/wp-sqlite
go run main.go
```

### scripts
- scripts under `script` directory are to be run from root directory e.g `bash script/test.sh` instead of `cd script; bash test.sh`
- docker build using script requires binary to built first, hence
1. `bash script/build.sh`
2. `bash script/dockerize.sh`
- running on docker not tried locally, some issue with my docker container network to host connection LOL but i would assume it works on machine where docker is properly set up

### problem assessment
- probably emulating chat application
- chat would use websocket for sending as well most likely, so slight modification to the usual flow

### structure 
an attempt to follow
- [Uncle Bob's Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)
- [Domain Driven Design](https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice)
- [Hexagonal](https://netflixtechblog.com/ready-for-changes-with-hexagonal-architecture-b315ec967749)

```bash
├── cmd                       	# entrypoint to application, allow for different implementations
│   └── wp-sqlite				# in this case, using SQLite for persistence
│   │   ├── main.go
│       └── ...
├── internal
│   ├── service					# service layer is for business logic
│   │   ├── message.go
│   │   └── message_test.go
│   ├── model                	# model layer is for shared model
│   │   └── message.go
│   ├── repository				# repository layer is for data persistence
│   │   └── sqlite				# in this case, using SQLite 
│   │       ├── message.go
│   │       └── message_test.go
│   └── handler					# handler layer is for interfaces with external application
│       ├── rest              	# in this case, using REST
│       │   ├── message.go
│       │  	└── message_test.go
│       └── websocket			# in this case, using websocket 
│           ├── message.go
│           └── message_test.go
├── pkg                       	# globally shared package
│       └── validator			# validator to validate request body 
│           ├── validate.go
│           ├── validate_test.go
│           └── ...
└── ...
```

### choice of tools
#### websocket/mqtt
- based on quick searching, mqtt is a layer above websocket, [go mqtt library](https://github.com/eclipse/paho.mqtt.golang) depends on [go websockets library](https://github.com/gorilla/websocket), hence using websocket in this exercise
- websocket looks simpler and easier to implement for now
#### persistence layer
- while usage of database is not required, getting all previously sent messages calls for some persistence layer. could have use cache but using embedded database instead

### misc
#### interface placement
- declaring interface in the place where it is consumed, as per [this talk](https://github.com/ronaudinho/iig)
#### short, repeated naming
- [only for short-lived variables](https://talks.golang.org/2014/names.slide#18)
#### testing
- black box testing preferance
- should have done white box testing for one unexported function in pkg/validator but nvm
