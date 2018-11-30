#!/usr/bin/env bash

operator-sdk build $1

echo "pushing $1"
docker push $1

kubectl delete -f ./deploy/crds/example_v1beta1_website_cr.yaml
kubectl create -f ./deploy/crds/example_v1beta1_website_cr.yaml

kubectl delete -f ./deploy/operator.yaml
kubectl create -f ./deploy/operator.yaml