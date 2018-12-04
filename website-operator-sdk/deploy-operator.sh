#!/usr/bin/env bash

deleteAndCreate=true
runLocal=false
operatorName="website-operator"
namespace="default"
image=architechbootcamp/website-operator:1.0.0

function log {
    echo "deploy-operator.sh --> $*"
}

function usage {
    log "To build, push, and deploy the operator to a cluster: ./deploy-operator.sh"
    log "To build, push, run operator locally: ./deploy-operator.sh -l"
}

while getopts ":l:" opt; do
    case $opt in
        l)
            runLocal=true
            log "will run operator locally"
            ;;
        \?) #invalid option
            log "${OPTARG} is not a valid option"
            usage
            exit 1
            ;;
    esac
done

if [ -z "$image" ] ; then log "docker image must be provided"; exit 1; fi

log "building the operator with image tag ${image}"
operator-sdk build "$image"

log "pushing operator image ${image}"
docker push "$image"

if [ "$deleteAndCreate" = true ]; then
    log "deleting the existing CRD"
    kubectl delete -f ./deploy/crds/example_v1beta1_website_crd.yaml
fi

log "creating Website CRD"
kubectl create -f ./deploy/crds/example_v1beta1_website_crd.yaml

if [ "$runLocal" = true ]; then
    log "running operator locally"
    export OPERATOR_NAME="$operatorName"
    operator-sdk up local --namespace "$default"
else
    log "deploying operator to the cluster"
    kubectl delete -f ./deploy/operator.yaml
    kubectl create -f ./deploy/operator.yaml
fi