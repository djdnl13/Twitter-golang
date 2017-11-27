#!/bin/bash

docker service create \
   --constraint node.role==manager \
   --replicas 1 --name dvizz -p 6969:6969 \
   --mount type=bind,source=/var/run/docker.sock,target=/var/run/docker.sock \
   --network my_network \
   eriklupander/dvizz
