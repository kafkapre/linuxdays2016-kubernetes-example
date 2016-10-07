#! /bin/bash

ids=$(docker ps -qf "name=k8s_server*")

IFS="\n" read -a idArr <<< "$ids"

docker rm -f ${idArr[0]}

curl $MINIKUBE_IP:30061/persons