# BARTENDER

this repository works in complement with: https://github.com/dicaormu/bartenderAsFunction

this is in charge of serving commands and keep the score

## execute the pipeline

aws cloudformation create-stack --stack-name bartender-deploy-pipeline --template-body file://pipeline/pipeline.yml --parameters ParameterKey=OAuthToken,ParameterValue=KEY_GITHUB  --capabilities CAPABILITY_IAM --profile PROFILE



 aws cloudformation delete-stack --role-arn arn:aws:iam::010154155802:role/bartender-deploy-pipeline-CFNRole-619RAE7YRDC7  --stack-name Bartender-Database --profile xebia
