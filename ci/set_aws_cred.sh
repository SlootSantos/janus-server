#!/usr/bin/env bash

mkdir ~/.aws 
AWS_CRED_FILE=~/.aws/credentials

echo "[default]" > $AWS_CRED_FILE
echo -e "aws_access_key_id=$INFRA_AWS_ACCESS_KEY" >> $AWS_CRED_FILE
echo -e "aws_secret_access_key=$INFRA_AWS_SECRET_KEY" >> $AWS_CRED_FILE

cat $AWS_CRED_FILE

