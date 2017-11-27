# Twitter-golang

## Requirements
This example were developed using:
 - go1.9.2
 - Jquery
 - BoltDB
 - Docker
 - Docker-machine
 - Docker-swarm
 - AMQP
 - Spring Cloud
 - Viper
 - Go Convey

## Setting up Docker Swarm cluster

    docker-machine create --driver virtualbox --virtualbox-cpu-count 2 --virtualbox-memory 2048 --virtualbox-disk-size 20000 swarm-manager-1
    eval "$(docker-machine env swarm-manager-1)"
    docker network create --driver overlay my_network
    docker swarm init --advertise-addr 192.168.99.100    
    
## Deploy cloud services
From /twitter

    go get
    ./keystore.sh
    ./dvizz.sh
    ./springcloud.sh
    ./support.sh

## Deploy microservices

    ./copyall.sh

## View running services

    docker service ls

## Demo
Demo files are in **/twitter/frontend**, you can use lampp or xampp or run it locally on your browser: **index.html** and **home.html**

## Author
This twitter microservice application was made at San Pablo Catholic University (Peru) for the Cloud Computing class by:
- [Daniel Lozano](https://github.com/djdnl13)
