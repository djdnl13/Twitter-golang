#!/bin/bash
export GOOS=linux
export CGO_ENABLED=0

cd accountservice;go get;go build -o accountservice-linux-amd64;echo built `pwd`;cd ..
cd tweetservice;go get;go build -o tweetservice-linux-amd64;echo built `pwd`;cd ..
cd healthchecker;go get;go build -o healthchecker-linux-amd64;echo built `pwd`;cd ..

cp healthchecker/healthchecker-linux-amd64 accountservice/
cp healthchecker/healthchecker-linux-amd64 tweetservice/

docker build -t djdnl13/accountservice accountservice/
docker build -t djdnl13/tweetservice tweetservice/

docker service rm accountservice
docker service rm tweetservice
docker service create --name=accountservice --replicas=1 --network=my_network -p=6767:6767 djdnl13/accountservice
docker service create --name=tweetservice --replicas=1 --network=my_network -p=6768:6768 djdnl13/tweetservice
