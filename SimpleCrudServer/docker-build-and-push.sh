#! /bin/bash

currentDir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

IMAGE_NAME="kafkapre/linuxdays2016-simple-crud-server:latest"

docker build -t $IMAGE_NAME  $currentDir

docker push $IMAGE_NAME
