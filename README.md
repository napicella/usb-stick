# USB stick
Store files and directories in the cloud to transfer it to another machine.   

The scope of the project is what I call a Sunday afternoon project, which means:
- few tests
- speed over code quality
- showcase technology, patterns, etc.
- useful to some extend
- whoever is interested should be able to build and deploy it with minimal effort
- maybe a live demo 

The project showcases:
- AWS Lambda with NodeJS which generates S3 presigned URL
- Cloudformation template to generate the usb service
- Cloudformation template to generate alarms
- a poor's man kill switch for API Gateway used for the demo
- golang CLI using [cobra](https://github.com/spf13/cobra)
- golang CLI dependencies vendored with [dep](https://github.com/golang/dep)

__Use cases__
- Copy data between two different machines which are in two different not connected VPCs
- Copy heapdump, logs, etc. from a prod machine down to your developer machine for deep investigations
- You use AWS Systems Manager Session Manager to ssh to your instances 

__what about scp (secure copy)?__  
__scp__ requires the source host and destination host to be routable, i.e. the hosts are in the same network or they have a public addressable ip.  
UsbStick stores the data in AWS S3, which serves as public addressable data buffer for the source and the destination host.
  
__what about Dropbox ?__  
Although you could use dropbox for this task, it comes with some drawbacks:
- The client is heavy. By comparison the usb-stick client is only few KB and it only requires zip to be installed on the 
machine.
- There is no official cli client for dropbox.
- Dropbox serves a different purpose: keep files in the cloud for as long as we want and there is no way to use
lifecycle policies on the data stored. On the other hand, with UsbStick we can enforce deletion of the data
after it was download or delete it after one hour from the date of the upload.

# Try the live demo

#### Install the client
__LINUX__

```bash
curl -L https://raw.githubusercontent.com/napicella/usb-stick/master/client_installer.sh | sudo bash -s linux
```

__MAC_OS__
```bash
curl -L https://raw.githubusercontent.com/napicella/usb-stick/master/client_installer.sh | sudo bash -s darwin
```

__WINDOWS__  
1. download the binary release/windows/usb.exe
2. set the executable flag
3. add it to the PATH


#### Usage
Store `some-folder` folder
```bash
usb upload ./some-folder
Enter Password:
```
This returns an etag. Something along the lines:
```bash
$ usb upload ./some-folder
Use this etag to download the file: 

2c8saf06-a064-12e9-a05b-d481d7a17d24
```

Download `some folder` in tmp
```bash
usb download -e 2c8saf06-a064-12e9-a05b-d481d7a17d24 -d /tmp
```

Storing and downloading individual files works exactly the same way.

# Deploy the stack to your account

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
The client is a golang cli which sends request to the Api Gateway.
The url of the Api Gateway is going to be injected in the goland code.
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
cd usb-stick-lambda-service && npm install
```

## Env variables

Set some variables for convenience

```bash
service_b=second.bucket.here
data_b=third.bucket.here
```

## Deploy the service

```bash
./deploy.sh -s "$service_b" -d "$data_b" -f
```

You can install the clients on how many machines as you want.   
Just remember the bucket to which you stored the client code :)

## Deploy alarms (optional but recommended)
Create an alarm which is triggered when too many requests are made too the api gw or the bucket which holds
the uploads grows too much in size.  

Set some variables for convenience  

```bash
code_b=fourth.bucket.here
email=some-email@some-domain.com
```

Then run:

```bash
./deploy_alarms.sh --codeBucket "${code_b}" --email "${email}"
```
