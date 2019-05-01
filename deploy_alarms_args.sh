#!/bin/bash

#!/bin/bash

# Define help function
function help(){
    echo "usb-stick-deploy-alarms - Deploy alarms for the usb-stick application";
    echo "Usage example:";
    echo "usb-stick-deploy-alarms (-f|--force) boolean (-c|--codeBucket) string]";
    echo "Options:";
    echo "-h or --help: Displays this information.";
    echo "-f or --force boolean: force reuse existing buckets. Required.";
    echo "-c or --codeBucket string: the bucket which will host the code for tha lambda function. Required.";
    echo "-e or --email string: email to notify when alarms trigger. Required.";
    exit 1;
}

# Declare vars. Flags initalizing to 0.
force=false

# Execute getopt
ARGS=$(getopt -o "hfc:e:" -l "help,force,codeBucket:,email:" -n "usb-stick-deploy-alarms" -- "$@");

#Bad arguments
if [ $? -ne 0 ];
then
    help;
fi

eval set -- "$ARGS";

while true; do
    case "$1" in
        -h|--help)
            shift;
            help;
            ;;
        -f|--force)
           shift;
            force=true;
            ;;
        -c|--codeBucket)
            shift;
            if [ -n "$1" ];
            then
                code_bucket="$1";
                shift;
            fi
            ;;
        -e|--email)
            shift;
            if [ -n "$1" ];
            then
                email="$1";
                shift;
            fi
            ;;

        --)
            shift;
            break;
            ;;
    esac
done

if [ -z "$code_bucket" ]
then
    echo "codeBucket is required";
    help;
fi

if [ -z "$email" ]
then
    echo "email is required";
    help;
fi