### Go NATS example


go mod init mod

$ go get -u 
$ go get -u=patch 
$ go get github.com/robteix/testmod@v1.0.1

- Install protoc-getn-go: https://github.com/golang/protobuf
- Generate proto: ```protoc --go_out=. *.proto```
- Build docker images: ```docker-compose build```
- Start docker compose: ```docker-compose --compatibility up```
