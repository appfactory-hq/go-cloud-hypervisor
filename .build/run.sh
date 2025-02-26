#!/bin/bash

cloud-hypervisor --version

kvm-ok

go version

echo "Go mod path: $GOPATH/pkg/mod"

ls -la /data

cd /src

make coverage
