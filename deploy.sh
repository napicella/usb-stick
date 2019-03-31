#!/bin/bash

set -e

stack_name="usbstick-presigned-url"
my_dir="$(dirname "$0")"
source "$my_dir/deploy_args.sh"

function create_bucket() {
    local bucket_name=$1
    local fail_if_exists=$2
    local bucket_exists=false

    if aws s3 ls "s3://$bucket_name" 2>&1 | grep -q 'An error occurred'
    then
        echo "bucket does not exit or permission is not there to view it."
        bucket_exists=false
    else
        bucket_exists=true

    fi

    if [[ "$bucket_exists" = false ]]; then
        aws s3 mb "s3://$bucket_name"
        return
    fi

    if [[ "$bucket_exists" = true && "$fail_if_exists" = false ]]; then
        printf "$bucket_name already exist. Reusing it because '-f' was specified \n"
        return
    fi

    printf "$bucket_name already exist. Rerun the program with -f if you wish to reuse existing buckets \n"
    exit -1
}

function deploy_create_cfn_stack() {
    sam package --template-file ./cloudformation/template.yaml \
        --s3-bucket "$serviceBucket"  \
        --output-template /tmp/out.yaml

    sam deploy --debug --template-file /tmp/out.yaml \
        --stack-name "$stack_name" \
        --parameter-overrides BucketName="$dataBucket" \
        --capabilities CAPABILITY_IAM

    aws cloudformation wait stack-create-complete --stack-name "$stack_name"

#    # get API endpoint
#    API_ENDPOINT=$(aws cloudformation describe-stacks --stack-name $stack_name --query 'Stacks[0].Outputs[0].OutputValue')
#
#    # remove quotes
#    API_ENDPOINT=$(sed -e 's/^"//' -e 's/"$//' <<< $API_ENDPOINT)

    apigwurl=$(aws cloudformation list-exports --query "Exports[?Name==\`UsbStickApiUrl\`].Value" --no-paginate --output text)
}

fail_if_exists=true
if [ "$force" = true ]; then
    fail_if_exists=false
fi

create_bucket "$clientBucket" "$fail_if_exists"
create_bucket "$serviceBucket" "$fail_if_exists"
create_bucket "$dataBucket" "$fail_if_exists"

deploy_create_cfn_stack


printf "Url of the service: \n$apigwurl \n"

"$my_dir"/build_client.sh -b "$clientBucket" -u "$apigwurl"
