#!/usr/bin/env bash
# Install helm if you have not already done so.
# See https://docs.helm.sh/using_helm/#installing-helm for your OS.

#1.  Add the coreos helm repo to get the Prometheus Operator chart
helm repo add coreos https://s3-eu-west-1.amazonaws.com/coreos-charts/stable/

#2. We will deploy everything to the monitoring namespace
kubectl create namespace monitoring

#3 Install the helm chart for the operator
helm install coreos/prometheus-operator --name prometheus-operator --namespace monitoring

#4. Install the actual Prometheus and Grafana pods.
helm install coreos/kube-prometheus --name kube-prometheus --set global.rbacEnable=true --namespace monitoring

#5. Create a Service to access the Grafana dashboard
kubectl create -f dashboard-service.yaml -n monitoring
