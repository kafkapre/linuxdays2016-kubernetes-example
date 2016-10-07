#! /bin/bash

currentDir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

source $currentDir/namespace.conf

echo "Creating namespace ..."
sed -e "s/%namespaceName%/$namespace/" $currentDir/namespace.yaml | kubectl create -f -


echo "Starting redis instance ..."
kubectl --namespace=$namespace create -f $currentDir/redis.yaml




case "$1" in
-s )
   echo "Starting SimpleCrudServer instance with service..."
   kubectl --namespace=$namespace create -f $currentDir/simpleCrudServer-with-service.yaml
   ;;
* )
   echo "Starting SimpleCrudServer instance without service..."
   kubectl --namespace=$namespace create -f $currentDir/simpleCrudServer.yaml
   ;;
esac


#echo "Starting SimpleCrudServer instance ..."
#kubectl --namespace=$namespace create -f $currentDir/simpleCrudServer-with-service.yaml
