#!/bin/bash

# Define help function
function help(){
    echo "usb-stick - download subcommand";
    echo "Usage example:";
    echo "usb-stick (-e|--etag) string (-d|--dir) string [(-h|--help)]";
    echo "Options:";
    echo "-h or --help: Displays this information.";
    echo "-e or --etag string: etag of the files. Required.";
    echo "-d or --dir string: download the files in the dir. Required.";
    exit 1;
}

# Declare vars. Flags initalizing to 0.

# Execute getopt
ARGS=$(getopt -o "he:d:" -l "help,etag:,dir:" -n "usb-stick" -- "$@");

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
        -e|--etag)
            shift;
                    if [ -n "$1" ];
                    then
                        etag="$1";
                        shift;
                    fi
            ;;
        -d|--dir)
            shift;
                    if [ -n "$1" ];
                    then
                        dir="$1";
                        shift;
                    fi
            ;;

        --)
            shift;
            break;
            ;;
    esac
done

# Check required arguments
if [ -z "$etag" ]
then
    echo "etag is required";
    help;
fi

if [ -z "$dir" ]
then
    echo "dir is required";
    help;
fi