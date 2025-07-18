FROM alpine
WORKDIR /app

ADD pkg pkg
ADD build build

ENTRYPOINT build/trip-service