FROM golang:alpine

RUN apk --no-cache add curl bash git

WORKDIR /go/src/autoscaler
ADD *.go ./
ADD deps.sh .
RUN ./deps.sh
RUN go install autoscaler

ENTRYPOINT autoscaler
