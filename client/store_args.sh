#!/bin/bash

# Define help function
function help(){
    echo "usb-stick - store - subcommand";
    echo "Usage example:";
    echo "usb-stick (-d|--dir) string (-p|--password) string [(-h|--help)]";
    echo "Options:";
    echo "-h or --help: Displays this information.";
    echo "-d or --dir string: path to the directory or to the file to store. Required.";
    echo "-p or --password string: used to protect the files. Required.";
    exit 1;
}

# Execute getopt
ARGS=$(getopt -o "hd:p:" -l "help,dir:,password:" -n "usb-stick" -- "$@");

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
        -d|--dir)
            shift;
                    if [ -n "$1" ];
                    then
                        dir="$1";
                        shift;
                    fi
            ;;
        -p|--password)
            shift;
                    if [ -n "$1" ];
                    then
                        password="$1";
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
if [ -z "$dir" ]
then
    echo "dir is required";
    help;
fi

if [ -z "$password" ]
then
    echo "password is required";
    help;
fi
