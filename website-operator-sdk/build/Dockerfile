FROM alpine:3.8

RUN apk upgrade --update --no-cache

USER nobody

ADD build/_output/bin/website-operator-sdk /usr/local/bin/website-operator-sdk
