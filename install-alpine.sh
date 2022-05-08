#!/bin/sh

set -e

mkdir /usr/local/share/ca-certificates
wget --no-check-certificate "https://ataboymirror-agct.gray.net/reliable_deployments/ssl_certs/rootSHA256.cer" -O "/usr/local/share/ca-certificates/rootSHA256.crt"
wget --no-check-certificate "https://ataboymirror-agct.gray.net/reliable_deployments/ssl_certs/issuingca1SHA256.cer" -O "/usr/local/share/ca-certificates/issuingca1SHA256.crt"
wget --no-check-certificate "https://ataboymirror-agct.gray.net/reliable_deployments/ssl_certs/issuingca2SHA256.cer" -O "/usr/local/share/ca-certificates/issuingca2SHA256.crt"
cd /usr/local/share/ca-certificates
cat rootSHA256.crt issuingca1SHA256.crt issuingca2SHA256.crt >> /etc/ssl/certs/ca-certificates.crt
cd -

echo Installing required packages
apk add --no-cache libstdc++ libaio ca-certificates bash tzdata unzip

export GLIBC_VERSION=2.26-r0

echo Installing glibc
wget --no-check-certificate "https://github.com/andyshinn/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-${GLIBC_VERSION}.apk" -O /tmp/glibc-${GLIBC_VERSION}.apk
wget --no-check-certificate "https://github.com/andyshinn/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-bin-${GLIBC_VERSION}.apk" -O /tmp/glibc-bin-${GLIBC_VERSION}.apk
wget --no-check-certificate "https://github.com/andyshinn/alpine-pkg-glibc/releases/download/${GLIBC_VERSION}/glibc-i18n-${GLIBC_VERSION}.apk" -O /tmp/glibc-i18n-${GLIBC_VERSION}.apk
apk add --allow-untrusted /tmp/*.apk
