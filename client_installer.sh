## usb-client installer for linux and mac
#!/usr/bin/env bash
OS=$1

if [[ -z "${OS}" ]]
then
	echo "Please set OS"
	exit -1
fi

wget "https://github.com/napicella/usb-stick/blob/master/release/${OS}/usb?raw=true" -O /usr/local/bin/usb
chmod +x /usr/local/bin/usb