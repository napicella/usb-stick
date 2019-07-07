#!/bin/bash

# Define help function
function help(){
    echo "usb-stick-deploy - Deploy the usb-stick application";
    echo "Usage example:";
    echo "usb-stick-deploy (-f|--force) boolean (-s|--serviceBucket) string (-d|--dataBucket) string [(-h|--help)]";
    echo "Options:";
    echo "-h or --help: Displays this information.";
    echo "-f or --force boolean: force reuse existing buckets. Required.";
    echo "-s or --serviceBucket string: the bucket which will host the service code. Required.";
    echo "-d or --dataBucket string: the bucket which will host the data. Required.";
    exit 1;
}

# Declare vars. Flags initalizing to 0.
force=false

# Execute getopt
ARGS=$(getopt -o "hfs:d:" -l "help,force,service-bucket:,data-bucket" -n "usb-stick-deploy" -- "$@");

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
        -s|--service-bucket)
            shift;
                    if [ -n "$1" ];
                    then
                        serviceBucket="$1";
                        shift;
                    fi
            ;;
        -d|--data-bucket)
            shift;
                    if [ -n "$1" ];
                    then
                        dataBucket="$1";
                        shift;
                    fi
            ;;

        --)
            shift;
            break;
            ;;
    esac
done


if [ -z "$serviceBucket" ]
then
    echo "service-bucket is required";
    help;
fi

if [ -z "$dataBucket" ]
then
    echo "data-bucket is required";
    help;
fi