FROM alpine
WORKDIR /app

ADD pkg pkg
ADD build build

ENTRYPOINT build/payment-service