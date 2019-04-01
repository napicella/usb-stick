# Usb stick
Store file and directories in the cloud to transfer it to another machine.  

__Use cases__
- Copy data between two different machines which are in two different not connected VPCs
- Copy heapdump, logs, etc. from a prod machine down to your developer machine for deep investigations 

__what about scp ?__  
__scp__ requires the source host and destination host to be routable, i.e. the hosts are in the same network or they have a public addressable ip.  
UsbStick stores the data in AWS S3, which serves as public addressable data buffer for the source and the destination host.
  
__what about dropbox ?__  
Although you could use dropbox for this task, it comes with some drawbacks:
- The client is heavy. By comparison the usb-stick client is only few KB and it only requires zip to be installed on the 
machine.
- There is no official cli client for dropbox.
- Dropbox serves a different purpose: keep files in the cloud for as long as we want and there is no way to use
lifecycle policies on the data stored. On the other hand, with UsbStick we can enforce deletion of the data
after it was download or delete it after one hour from the date of the upload.

# Usage
You need to follow the __Getting Started__ section to create the service and the client.  
Assuming that, here it's how you would use it:  
Install the client on a machine:

```bash
wget "http://your_code_client_bucket.s3.amazonaws.com/usb-stick.zip" -O /tmp/usb-stick.zip; \
      sudo unzip /tmp/usb-stick.zip -d /usr/bin/; \
      rm /tmp/usb-stick.zip
```
This downloads the client from the bucket it was uploaded during the deployment (as per __Getting Started__ section ) and 
installs the __usbstick__ program under __/usr/bin__

Store a folder
```bash
usbstick store -d ./some-folder -p some-password-used-to-encrypt-the-zip
```
This returns an etag. Something along the lines:
```bash
$ usbstick store -d ./some-folder -p 1234
  adding: some-folder (deflated 81%)
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:--  0:00:01 --:--:--     0
Use this etag to download the file: 

Em5lXDayaslXdrVwjKokETzU-eu26tCr
```

Download `some folder` in tmp
```bash
usbstick download -e Em5lXDayaslXdrVwjKokETzU-eu26tCr -d /tmp
```


# Getting started

The guide describes how to deploy the UsbStick Service and Client in your AWS account.
__Prerequisites__  
- Aws cli installed on your machine
- Sam installed on your machine
- npm installed on your machine

### What are we going to create? 
__Service__
The cloudformation template in the source creates the service:
- An Api Gateway which proxy requests to a Lambda function.  
- The Lambda function creates S3 Presigned url to upload and download objects from the bucket we want to use as storage for 
our cloud Usb Stick.

__Client__  
The client is a bash script which sends request to the Api Gateway.
The url of the Api Gateway is going to be injected in the script and the script is going to be uploaded to an S3 bucket. 
During upload the client performs:
- creates a random etag which is going to be used as name of the object in S3
- zip the content to upload using a password, send http requests to the API Gateway and uses the presigned url to upload the zip.
- The etag gets printed on the screen, so it can be used as parameter for the download command

During download:
- Send an http request to the Api Gateway with the etag 
- Client gets the presigned url and start the download.
- Upon completion, unzip the content


## Build lambda code
```
cd src && npm install
```

## Env variables

Set some variables for convenience

```bash
client_b=first.bucket.here
service_b=second.bucket.here
data_b=third.bucket.here
```

## Deploy the service and the client

```bash
./deploy.sh -c "$client_b" -s "$service_b" -d "$data_b" -f
```

## Install the client

```bash
wget "http://$client_b.s3.amazonaws.com/usb-stick.zip" -O /tmp/usb-stick.zip; \
      sudo unzip /tmp/usb-stick.zip -d /usr/bin/; \
      rm /tmp/usb-stick.zip
```

You can install the clients on how many machines as you want.   
Just remember the bucket to which you stored the client code :)