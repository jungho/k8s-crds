#!/usr/bin/env bash

operator-sdk build $1

echo "pushing $1"
docker push $1