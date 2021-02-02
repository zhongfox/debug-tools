# GRPC-GO-DOCKER-HELLOWORLD


for from https://github.com/juaruipav/grpc-go-docker-helloworld

Run the server:
```
$ cd server 
$ docker build . -t zhongfox/grpc-server:v1
$ docker run --net=host -it zhongfox/grpc-server:v1
```


Run the client:
```
$ cd client 
$ docker build . -t zhongfox/grpc-client:v1
$ docker run --net=host -it zhongfox/grpc-client:v1
```

