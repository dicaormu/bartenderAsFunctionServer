# BARTENDER

this repository works in complement with: https://github.com/dicaormu/bartenderAsFunction

this is in charge of serving commands and keep the score

## execute the pipeline

aws cloudformation create-stack --stack-name bartender-deploy-pipeline --template-body file://pipeline/pipeline.yml --parameters ParameterKey=OAuthToken,ParameterValue=KEY_GITHUB  --capabilities CAPABILITY_IAM --profile PROFILE


