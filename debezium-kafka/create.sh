#! /bin/bash

curl -s -X POST -H "Accept:application/json" \
    -H "Content-Type:application/json" \
    localhost:8083/connectors/ -d @mysql.json | jq