#!/bin/bash

set -e

# Generate a default CA and Certificate

TIMESTAMP=$(date +%s)
CERTNAME="certificate-$TIMESTAMP"
TMPFILE="/tmp/credhub_test_cert_$TIMESTAMP.json"
CERTFILE="./tmp/$CERTNAME.cer"

mkdir -p tmp

./build/credhub generate -t certificate -n "ca" --common-name my-ca --is-ca --self-sign 2>&1 > /dev/null
./build/credhub generate -n $CERTNAME -t certificate --common-name $CERTNAME --ca "ca" -g digital_signature -g key_agreement -e code_signing -e email_protection 2>&1 > /dev/null
./build/credhub get -n $CERTNAME --output-json > $TMPFILE

printf "$(cat $TMPFILE | jq '.certificate')" | sed -e 's/"//g' >  $CERTFILE
printf "$(cat $TMPFILE | jq '.private_key')" | sed -e 's/"//g' >> $CERTFILE

echo "Generated $CERTNAME in $CERTFILE"
