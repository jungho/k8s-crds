#!/usr/bin/env bash

deleteAndCreate=true
runLocal=false
operatorName="website-operator"
namespace="default"

function usage {
    log "To build, push, and deploy the operator to a cluster: ./deploy_operator.sh <image>"
    log "To build, push, run operator locally: ./deploy_operator.sh -l <image>"
}

while getopts ":l" opt; do
    case $opt in
        l)
            runLocal=true
            ;;
        \?) #invalid option
            log "${OPTARG} is not a valid option"
            usage
            exit 1
            ;;
    esac
done

if [ -z "$1" ] ; then log "docker image must be provided"; exit 1; fi

echo "building the operator..."
operator-sdk build "$1"

echo "pushing operator image $1"
docker push "$1"

if [ "$deleteAndCreate" = true ]; then
    echo "deleting the existing CRD"
    kubectl delete -f ./deploy/crds/example_v1beta1_website_crd.yaml
fi

echo "creating Website CRD"
kubectl create -f ./deploy/crds/example_v1beta1_website_crd.yaml

if [ "$runLocal" = true ]; then
    echo "running operator locally"
    export OPERATOR_NAME="$operatorName"
    operator-sdk local --namespace "$default"
else
    echo "deploying operator to the cluster"
    kubectl delete -f ./deploy/operator.yaml
    kubectl create -f ./deploy/operator.yaml
fi