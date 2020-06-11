#!/bin/bash

# sourcing nvm
. ~/.nvm/nvm.sh

ENV_MAP="{\"dev_env\":\"dev\",\"stage_env\":\"stage\",\"prod_env\":\"prod\"}"

echo "BRANCH IS: $STACKERS_BRANCH"
echo "COMMIT IS: $STACKERS_COMMIT"
echo "PR_ID IS: $STACKERS_PR_ID"


echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Cloning Repository \x1B[0m"
git clone https://${OAUTH_TOKEN}:x-oauth-basic@github.com/${REPO_FULL}.git

cd ${REPO}

if [ "$STACKERS_COMMIT" != "" ]
then
  git checkout $STACKERS_COMMIT
fi

echo -e "\x1B[30;48;5;82m Cloning Repository \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"

OUTPUT_DIR="build/"

find .stackers/config.json
if [ $? -eq 0 ]
then
  echo "FOUND CONFIG"
  echo "Setting Node Version"
  CONFIG_NODE_VERSION=$(cat .stackers/config.json | jq -r '.node_version')
  if [ "$CONFIG_NODE_VERSION" == "null" ]
  then 
    CONFIG_NODE_VERSION="lts/*"
  fi

  nvm install $CONFIG_NODE_VERSION
  if [ $? -ne 0 ]
  then
    echo "Using invalid Node version. Please choose one of the following:"
    nvm ls-remote

    echo -e "\n \n\x1B[48;5;9m                                    \x1B[97mDEPLOYMENT FAILED!                                    \x1B[0m"
    exit 1
  fi

  DEPLOY_BRANCH=$(cat .stackers/config.json | jq -r --arg BRANCH "${STACKERS_BRANCH}" '.branches | to_entries[] | select(.value == $BRANCH) | .key')

  echo SETTING ENV\'s
  eval $(cat .stackers/config.json | jq -r '.envs[] | to_entries[] | [.key,.value] | "export " + join("=")')

  echo "pre_install"
  eval $(cat .stackers/config.json | jq -r '.pre_install')

  echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Installing NPM packages \x1B[0m"
  npm install
  echo -e "\x1B[30;48;5;82m Installing NPM packages \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"

  echo "post_install"
  eval $(cat .stackers/config.json | jq -r '.post_install')

  echo "build"
  eval $(cat .stackers/config.json | jq -r '.build')

  OUTPUT_DIR=$(cat .stackers/config.json | jq -r '.output_dir')
else
  echo "NO CONFIG"

  echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Installing NPM packages \x1B[0m"
  npm install
  echo -e "\x1B[30;48;5;82m Installing NPM packages \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"

  echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Building Production \x1B[0m"
  npm run build 
  echo -e "\x1B[30;48;5;82m Building Production \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"
fi

DEPLOY_TARGET_ENV=$(echo $ENV_MAP | jq -r --arg DEPLOY_BRANCH "${DEPLOY_BRANCH}" 'to_entries[] | select(.key == $DEPLOY_BRANCH) | .value' )
if [ "$DEPLOY_TARGET_ENV" == "" ]
then
  if [ "$STACKERS_PR_ID" != "" ]
  then
    DEPLOY_TARGET_ENV="pr/$STACKERS_PR_ID"
  else
  echo "No deployment conifgured for branch: $STACKERS_BRANCH" 
  echo -e "\n \n\x1B[30;48;5;82m                                    DEPLOYMENT SUCCESSFUL!                                    \x1B[0m"
  exit 0
  fi
fi

# get new version for app
GREEN_VERSION=$(cat package.json | jq -r '.version')

echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Setting next app version \x1B[0m $GREEN_VERSION"
# upload new version to bucket
echo -e "{\"version\":\"$GREEN_VERSION\"}" | aws s3 cp - s3://$BUCKET/green/version.txt 

# upload build to green
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Uploading GREEN deployment \x1B[0m"
aws s3 sync ./$OUTPUT_DIR s3://$BUCKET/green > /dev/null
echo -e "\x1B[30;48;5;82m Uploading GREEN deployment \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"


# invalidate CDN for /green
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Invalidating CDN \x1B[0m"
aws cloudfront create-invalidation --distribution-id=${CDN} --paths=/* > /dev/null
echo -e "\x1B[30;48;5;82m Invalidating CDN \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"

# wait for invalidation
sleep 10

## Get cloudfront domain (or pass it?)
CNAME=$(aws cloudfront get-distribution-config --id=$CDN | jq -r '.DistributionConfig.Aliases.Items[] | select(test("green."))')

# ############### PART III ###################
#         ####### RUN INTEGRATION #######
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Run Integration test \x1B[0m"
JSON_OUTPUT=$(node /home/pipeline/util/integration.js https://$CNAME | jq -R 'fromjson?')
ERROR_OUTPUT=$(echo $JSON_OUTPUT | jq -r '.error')

echo -e "Runtime output: \n $JSON_OUTPUT"

# echo "OUPUT______________________ $JSON_OUTPUT"
# rollback if broken app
if [ "$ERROR_OUTPUT" != "" ] && [ "$ERROR_OUTPUT" != "null" ]; then
  echo -e "\x1B[31mBroken GREEN deployment. Not swapping blue <> green. \x1B[0m"
  echo -e "\x1B[4mRuntime Error:\x1B[24m $ERROR_OUTPUT"
  echo -e "\n \x1B[36mYou can reproduce the error @ \x1B[0m https://$CNAME"

  echo -e "\n \n\x1B[48;5;9m                                    \x1B[97mDEPLOYMENT FAILED!                                    \x1B[0m"


  exit 1
fi

## curl distro.domain/version.txt (with green prefix)
DEPLOYED_VERSION=$(curl https://$CNAME/version.txt | jq -r '.version')
## if version == curr_version => success (or only print response for now...)
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Checking deployed version: \x1B[0m $DEPLOYED_VERSION"
# rollback if not correct version
# maybe do retry
if [ "$DEPLOYED_VERSION" != "$GREEN_VERSION" ]; then
  echo "not the same version => rollback"

  exit 1
fi
echo -e "\x1B[30;48;5;82m Checking deployed version \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"


aws s3 ls s3://$BUCKET/$DEPLOY_TARGET_ENV/latest/ 
if [[ $? -eq 0 ]]; then
    CURR_VERSION=$(aws s3 cp s3://$BUCKET/$DEPLOY_TARGET_ENV/latest/version.txt - | jq -r '.version') # what if no version.json in bucket?

    echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Storing 'Latest' as v$CURR_VERSION \x1B[0m"
    aws s3 mv s3://$BUCKET/$DEPLOY_TARGET_ENV/latest s3://$BUCKET/$DEPLOY_TARGET_ENV/v$CURR_VERSION --recursive  > /dev/null
    echo -e "\x1B[30;48;5;82m Storing 'Latest' as v$CURR_VERSION \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"
fi

echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Deploying GREEN deployment to 'Latest' \x1B[0m"
aws s3 mv s3://$BUCKET/green s3://$BUCKET/$DEPLOY_TARGET_ENV/latest --recursive > /dev/null
echo -e "\x1B[30;48;5;82m Deploying GREEN deployment to 'Latest \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"


# finally invalidate cdn
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Invalidating CDN \x1B[0m"
aws cloudfront create-invalidation --distribution-id=${CDN} --paths=/* > /dev/null
echo -e "\x1B[30;48;5;82m Invalidating CDN \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"



echo -e "\n \n\x1B[30;48;5;82m                                    DEPLOYMENT SUCCESSFUL!                                    \x1B[0m"


