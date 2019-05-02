#!/bin/bash

set -e

# Define help function
function help(){
    echo "usb-stick-client-build - Build the usb-stick client";
    echo "Usage example:";
    echo "usb-stick-client-build (-b|--bucket) string (-u|--url) string [(-h|--help)]";
    echo "Options:";
    echo "-h or --help: Displays this information.";
    echo "-b or --bucket string: where to upload the client. Required.";
    echo "-u or --url string: url of the usb stick service. Required.";
    exit 1;
}

# Declare vars. Flags initalizing to 0.

# Execute getopt
ARGS=$(getopt -o "hb:u:" -l "help,bucket:,url:" -n "usb-stick-client-build" -- "$@");

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
        -b|--bucket)
            shift;
                    if [ -n "$1" ];
                    then
                        bucket="$1";
                        shift;
                    fi
            ;;
        -u|--url)
            shift;
                    if [ -n "$1" ];
                    then
                        url="$1";
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
if [ -z "$bucket" ]
then
    echo "bucket is required";
    help;
fi

if [ -z "$url" ]
then
    echo "url is required";
    help;
fi


mkdir -p release
cd client
echo "$(pwd)"
echo "$bucket"
echo "$url"

sed -i -e "s~API_URL~$url~" usbstick
zip -r ../release/usb-stick.zip *
sed -i -e "s~$url~API_URL~" usbstick
aws s3 cp ../release/usb-stick.zip "s3://$bucket/usb-stick.zip"
aws s3api put-object-acl --bucket "$bucket" --key usb-stick.zip --acl public-read

cat >../release/installer.sh <<EOL
#!/bin/bash
wget "http://${bucket}.s3.amazonaws.com/usb-stick.zip" -O /tmp/usb-stick.zip;
sudo unzip -o /tmp/usb-stick.zip -d /usr/local/bin/;
rm /tmp/usb-stick.zip
EOL
chmod +x ../release/installer.sh
