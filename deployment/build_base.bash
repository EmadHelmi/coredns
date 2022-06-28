#!/bin/bash
docker build --build-arg HTTP_PROXY=http://185.244.213.113:1398 --build-arg HTTPS_PROXY=http://185.244.213.113:1398 -t coredns_base:latest -f BaseDockerfile .
