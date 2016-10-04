#! /bin/bash

currentDir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

source $currentDir/namespace.conf

namespace=linuxdays

echo "Creating namespace ..."
sed -e "s/%namespaceName%/$namespace/" $currentDir/namespace.yaml | kubectl create -f -


echo "Starting redis instance ..."
/opt/kubectl --namespace=$namespace create -f $currentDir/redis.yaml
