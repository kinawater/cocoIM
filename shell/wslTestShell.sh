#!/bin/bash

docker run -v $(pwd)/../cocoIM:/cocoIM -v $(pwd)/../config:/config --name IMserviceTest --privileged --cap-add=NET_ADMIN -d alpine /cocoIM gateway --config="./config/config.yaml"


