#!/usr/bin/env bash

# get parameter from AWS parameter store
getSSMParam() {
    echo $1
    PARAM_VAL=$(aws ssm get-parameter --name janus-env-$1 --region us-east-1 | jq -r '.Parameter.Value')
    echo "$1=$PARAM_VAL" >> .new_env
}

# split eachline at the "=" so all env's get overwritten during script
getValueForEnv() {
    IFS='=' 
    read -ra ENV_VAR <<< "$1"
    getSSMParam $ENV_VAR
}

# read every line from .env.example file to get each variable name
while read line;
    do
        if [ -z "$line" ]; then
            continue
        fi

        getValueForEnv $line;
    done < $PWD/.env
