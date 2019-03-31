# Usb stick
Store file and directories in the cloud to transfer it to another machine.
__what about scp ?__
__what about dropbox ?__

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

You can install the clients on how many machine as you want.   
Just remember the bucket to which you stored the client code :)