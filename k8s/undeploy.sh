#! /bin/bash

currentDir=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )

source $currentDir/namespace.conf

echo "deleting RCs and SVCs ..."
kubectl --namespace=$namespace delete replicationcontrollers,pods,services --grace-period=0 -l "app=crud-server"

echo "Deleting namespace ..."
kubectl delete namespaces $namespace
