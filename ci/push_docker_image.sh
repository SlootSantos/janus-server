#!/usr/bin/env bash

# login to ECR
$(aws ecr get-login --no-include-email --region us-east-1)

docker build -t janus/server:latest -f "${PWD}/Dockerfile" ${PWD}

IMAGE_ID=$(docker images -q janus/server:latest)

docker images ls

docker tag $IMAGE_ID 108151951856.dkr.ecr.us-east-1.amazonaws.com/janus/server:latest
docker push 108151951856.dkr.ecr.us-east-1.amazonaws.com/janus/server:latest
