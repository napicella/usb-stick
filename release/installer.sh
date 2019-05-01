#!/bin/bash
wget "http://com.napicella.usbstick.client.code.s3.amazonaws.com/usb-stick.zip" -O /tmp/usb-stick.zip;
sudo unzip -o /tmp/usb-stick.zip -d /usr/bin/;
rm /tmp/usb-stick.zip
