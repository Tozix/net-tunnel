#! /bin/sh
GOOS=linux GOARCH=mipsle GOMIPS=softfloat go build -a -ldflags "-s -w" -o $2.tip "$1"
upx -f --brute -o $2 $2.tip
