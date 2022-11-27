FROM golang:alpine3.16 as builder

RUN apk add --no-cache --virtual .build-deps gcc musl-dev openssl git

RUN mkdir /go/src/github.com
RUN mkdir /go/src/github.com/cheetahfox

WORKDIR /go/src/github.com/cheetahfox

RUN git clone https://github.com/cheetahfox/singularity-iot.git

WORKDIR /go/src/github.com/cheetahfox/singularity-iot/

RUN go build

FROM alpine:3.16

COPY --from=builder /go/src/github.com/cheetahfox/singularity-iot/singularity-iot . 
EXPOSE 2200
CMD ./singularity-iot