#!/usr/bin/env bash

ENV_FILE_INPUT=$PWD/.env.example
ENV_FILE_OUTPUT=$PWD/.env

# get parameter from AWS parameter store
getSSMParam() {
    echo ---------------
    echo "Setting ENV => $1"
    PARAM_VAL=$(aws ssm get-parameter --name /janus/env/production/$1 --region us-east-1 --with-decryption | jq -r '.Parameter.Value')
    echo "$1=$PARAM_VAL" >> $ENV_FILE_OUTPUT
    echo -e "Done \xE2\x9C\x94"
}

# split eachline at the "=" so all env's get overwritten during script
getValueForEnv() {
    IFS='=' 
    read -ra ENV_VAR <<< "$1"
    getSSMParam $ENV_VAR
}

# delete if exisiting output file
rm $ENV_FILE_OUTPUT 2> /dev/null
# read every line from .env.example file to get each variable name
while read line;
    do
        if [ -z "$line" ]; then
            continue
        fi

        getValueForEnv $line;
    done < $ENV_FILE_INPUT
