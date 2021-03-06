#!/bin/bash

export GOOS=linux
export CGO_ENABLED=0

cd gelftail;go get;go build -o gelftail-linux-amd64;echo built `pwd`;cd ..

#export GOOS=darwin

docker build -t djdnl13/gelftail gelftail/
docker service rm gelftail
docker service create --name=gelftail -p=12202:12202/udp --replicas=1 --network=my_network djdnl13/gelftail
