#!/usr/bin/env bash

operator-sdk build $1

echo "pushing $1"
docker push $1

kubectl create -f ./deploy/service_account.yaml
kubectl create -f ./deploy/role.yaml
kubectl create -f ./deploy/role_binding.yaml
kubectl create -f ./deploy/operator.yaml