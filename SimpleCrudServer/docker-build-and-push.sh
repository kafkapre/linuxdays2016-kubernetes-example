#! /bin/bash

currentDir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

IMAGE_NAME="kafkapre/linuxdays2016-simple-crud-server"

IMAGE_ID=`docker build -t $IMAGE_NAME  $currentDir/`

docker tag $IMAGE_ID "$IMAGE_NAME:latest"

docker push $IMAGE_NAME
