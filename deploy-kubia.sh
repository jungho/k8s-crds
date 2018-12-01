#!/usr/bin/env bash

kubectl create -f website-crd.yaml
#The website-controller will need a serviceAccount with sufficient privileges to access the kube-apiserver
kubectl create serviceaccount website-controller
kubectl create clusterrolebinding website-controller --clusterrole=cluster-admin --serviceaccount=default:website-controller
#create the website-controller as a deployment
kubectl create -f website-controller.yaml
kubectl create -f website.yaml
