#!/bin/bash

git clone https://${OAUTH_TOKEN}:x-oauth-basic@github.com/${REPO_FULL}.git
cd ${REPO}

echo "Start Building"
npm install
npm run build 
ls
rm -rf node_modules/

aws s3 sync ./build/ s3://${BUCKET}
aws cloudfront create-invalidation --distribution-id=${CDN} --paths=/*