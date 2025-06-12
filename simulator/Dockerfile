FROM alpine:latest

RUN apk update && apk add freeradius-utils

WORKDIR /simulator

COPY acct_start.txt ./
COPY acct_stop.txt ./
