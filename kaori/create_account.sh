#!/bin/sh
if [ "$#" -ne 2 ]
then
	echo "Usage: $0 <username> <password>"
	exit 1
fi
curl\
	-X POST\
	-H Content-Type:application/json\
	-k\
	https://localhost:8080/account/create\
	-d "{\"Username\":\"$1\", \"Password\":\"$2\"}"

