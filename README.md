### Go NATS example


go mod init mod

$ go get -u 
$ go get -u=patch 
$ go get github.com/robteix/testmod@v1.0.1

- Install protoc-getn-go: https://github.com/golang/protobuf
- Generate proto: ```protoc --go_out=. *.proto```
- Build docker images: ```docker-compose build```
- Start docker compose: ```docker-compose --compatibility up```
- Docker stats: ```docker stats -a go-nats_jvmdealer_1 go-nats_dealer_1 go-nats_factory_1```
- Example of memory limit in docker-compose.yml: 
<pre>
factory:
  build: ./factory
  network_mode: host
  deploy:
    resources:
      limits:
        memory: 1024M
</pre>
