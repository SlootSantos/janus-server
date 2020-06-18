#!/usr/bin/env bash

mkdir ~/.aws 
AWS_CRED_FILE=~/.aws/credentials
echo "CIRCLE TAG & BRANCH: $CIRCLE_TAG $CIRCLE_BRANCH"
[[ "$CIRCLE_TAG" != "" ]] && ACCESS="$INFRA_AWS_ACCESS_KEY" || ACCESS="$INFRA_AWS_ACCESS_KEY_DEVELOP"
[[ "$CIRCLE_TAG" != "" ]] && SECRET="$INFRA_AWS_SECRET_KEY" || SECRET="$INFRA_AWS_SECRET_KEY_DEVELOP"

echo "[default]" > $AWS_CRED_FILE
echo -e "aws_access_key_id=$ACCESS" >> $AWS_CRED_FILE
echo -e "aws_secret_access_key=$SECRET" >> $AWS_CRED_FILE

cat $AWS_CRED_FILE

