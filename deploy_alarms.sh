#!/bin/bash

set -e

stack_name="usbstick-alarms"
my_dir="$(dirname "$0")"
source "$my_dir/deploy_alarms_args.sh"

bucket_to_monitor="UploadBucketName"
api_gateway_name="UsbStickApi"
api_gateway_id="UsbStickApiId"
api_gateway_stage="StageName"

if ! aws cloudformation describe-stacks --stack-name "$stack_name" ; then
    wait_condition="stack-create-complete"
else
    wait_condition="stack-update-complete"
fi

sam package --template-file ./cloudformation/alarms_template.yml \
    --s3-bucket "${code_bucket}" \
    --output-template /tmp/out.yaml
sam deploy --debug --template /tmp/out.yaml \
    --stack-name "${stack_name}" \
    --parameter-overrides \
    Email=napicella4@gmail.com ApiGatewayApiName="${api_gateway_name}" ApiGatewayStage="${api_gateway_stage}" ApiGatewayId="${api_gateway_id}" S3DataBucketName="${bucket_to_monitor}" \
    --capabilities CAPABILITY_IAM CAPABILITY_AUTO_EXPAND \
    --no-fail-on-empty-changeset

aws cloudformation wait "$wait_condition" --stack-name "$stack_name"
