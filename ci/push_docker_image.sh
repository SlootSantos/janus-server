#!/usr/bin/env bash


echo $REDIS_PORT
echo $REDIS_HOSTNAME
echo $PWD


# login to ECR
$(AWS_SECRET_ACCESS_KEY=$AWS_SECRET_KEY AWS_ACCESS_KEY_ID=$AWS_ACCESS_KEY_ID aws ecr get-login --no-include-email --region us-east-1)

docker build -t janus/server:latest -f "${PWD}/Dockerfile" ${PWD}

IMAGE_ID=$(docker images -q janus/server:latest)

docker tag $IMAGE_ID 108151951856.dkr.ecr.us-east-1.amazonaws.com/janus/server:latest
docker push 108151951856.dkr.ecr.us-east-1.amazonaws.com/janus/server:latest
