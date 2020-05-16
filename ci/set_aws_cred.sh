AWS_CRED_FILE=~/.aws/credentials

echo "[default]" > $AWS_CRED_FILE
echo "aws_access_key_id=$AWS_ACCESS_KEY_ID" >> $AWS_CRED_FILE
echo "aws_secret_access_key=$AWS_SECRET_ACCESS_KEY" >> $AWS_CRED_FILE

cat $AWS_CRED_FILE