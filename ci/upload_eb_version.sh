#!/usr/bin/env bash

ACCOUNT_ID=$(aws sts get-caller-identity | jq -r '.Account')

zip -r app_v_$CIRCLE_BUILD_NUM.zip .ebextensions/ Dockerrun.aws.json
aws s3 cp app_v_$CIRCLE_BUILD_NUM.zip s3://elasticbeanstalk-us-east-1-$ACCOUNT_ID/
aws elasticbeanstalk create-application-version --application-name janus_server --version-label v$CIRCLE_BUILD_NUM --description="New Version number $CIRCLE_BUILD_NUM" --source-bundle S3Bucket="elasticbeanstalk-us-east-1-$ACCOUNT_ID",S3Key="app_v_$CIRCLE_BUILD_NUM.zip" --auto-create-application --region=us-east-1
aws elasticbeanstalk update-environment --application-name janus_server --environment-name=janus-server-production --version-label v$CIRCLE_BUILD_NUM --region=us-east-1