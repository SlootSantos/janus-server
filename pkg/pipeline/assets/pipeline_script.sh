#!/bin/bash
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Cloning Repository \x1B[0m"
git clone https://${OAUTH_TOKEN}:x-oauth-basic@github.com/${REPO_FULL}.git
cd ${REPO}
echo -e "\x1B[30;48;5;82m Cloning Repository \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"

echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Installing NPM packages \x1B[0m"
npm install
echo -e "\x1B[30;48;5;82m Installing NPM packages \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"

echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Building Production \x1B[0m"
npm run build 
echo -e "\x1B[30;48;5;82m Building Production \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"


# get new version for app
GREEN_VERSION=$(cat package.json | jq -r '.version')

echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Setting next app version \x1B[0m $GREEN_VERSION"
# upload new version to bucket
echo -e "{\"version\":\"$GREEN_VERSION\"}" | aws s3 cp - s3://$BUCKET/green/version.txt 

# upload build to green
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Uploading GREEN deployment \x1B[0m"
aws s3 sync ./build s3://$BUCKET/green > /dev/null
echo -e "\x1B[30;48;5;82m Uploading GREEN deployment \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"


# invalidate CDN for /green
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Invalidating CDN \x1B[0m"
aws cloudfront create-invalidation --distribution-id=${CDN} --paths=/* > /dev/null
echo -e "\x1B[30;48;5;82m Invalidating CDN \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"

# wait for invalidation
sleep 10

## Get cloudfront domain (or pass it?)
CNAME=$(aws cloudfront get-distribution-config --id=$CDN | jq -r '.DistributionConfig.Aliases.Items[] | select(test("green-"))')

# ############### PART III ###################
#         ####### RUN INTEGRATION #######
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Run Integration test \x1B[0m"
JSON_OUTPUT=$(node /util/integration.js https://$CNAME | jq -R 'fromjson?')
ERROR_OUTPUT=$(echo $JSON_OUTPUT | jq -r '.error')

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


aws s3 ls s3://$BUCKET/latest/ 
if [[ $? -eq 0 ]]; then
    CURR_VERSION=$(aws s3 cp s3://$BUCKET/latest/version.txt - | jq -r '.version') # what if no version.json in bucket?

    echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Storing 'Latest' as v$CURR_VERSION \x1B[0m"
    aws s3 mv s3://$BUCKET/latest s3://$BUCKET/v$CURR_VERSION --recursive  > /dev/null
    echo -e "\x1B[30;48;5;82m Storing 'Latest' as v$CURR_VERSION \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"
fi

echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Deploying GREEN deployment to 'Latest' \x1B[0m"
aws s3 mv s3://$BUCKET/green s3://$BUCKET/latest --recursive > /dev/null
echo -e "\x1B[30;48;5;82m Deploying GREEN deployment to 'Latest \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"


# finally invalidate cdn
echo -e "\n \n\x1B[40;38;5;82m Step \x1B[30;48;5;82m Invalidating CDN \x1B[0m"
aws cloudfront create-invalidation --distribution-id=${CDN} --paths=/* > /dev/null
echo -e "\x1B[30;48;5;82m Invalidating CDN \x1B[0m           \x1B[38;5;82m\xE2\x9C\x94\x1B[0m"



echo -e "\n \n\x1B[30;48;5;82m                                    DEPLOYMENT SUCCESSFUL!                                    \x1B[0m"


