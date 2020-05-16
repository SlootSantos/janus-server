#!/usr/bin/env bash

# login to ECR
$(aws ecr get-login --no-include-email --region us-east-1)

docker build -t janus/server:latest -f "${PWD}/Dockerfile" ${PWD}

IMAGE_ID=$(docker images -q janus/server:latest)
ECR_URL=$(aws ssm get-parameter --name /janus/env/production/ECR_URL --region us-east-1 --with-decryption | jq -r '.Parameter.Value')

docker tag $IMAGE_ID $ECR_URL/janus/server:latest
docker push $ECR_URL/janus/server:latest
