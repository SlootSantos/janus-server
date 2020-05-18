#!/usr/bin/env bash

ACCOUNT_ID=$(aws sts get-caller-identity | jq -r '.Account')

aws s3 cp Dockerrun.aws.json s3://elasticbeanstalk-us-east-1-$ACCOUNT_ID/Dockerrun.aws.v_$CIRCLE_BUILD_NUM.json
aws elasticbeanstalk create-application-version --application-name janus_server --version-label v$CIRCLE_BUILD_NUM --description="New Version number $CIRCLE_BUILD_NUM" --source-bundle S3Bucket="elasticbeanstalk-us-east-1-$ACCOUNT_ID",S3Key="Dockerrun.aws.v_$CIRCLE_BUILD_NUM.json" --auto-create-application --region=us-east-1 --process
aws elasticbeanstalk update-environment --application-name janus_server --environment-name=janus-server-production --version-label v$CIRCLE_BUILD_NUM --region=us-east-1