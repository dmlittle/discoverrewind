FROM golang:1.16.0 as build

WORKDIR /go/src/github.com/dmlittle/discoverrewind

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o build/serve cmd/serve/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o build/migrations cmd/migrations/*.go

FROM alpine:3.12

RUN apk update \
  && apk upgrade \
  && apk add --no-cache ca-certificates \
  && update-ca-certificates

# Add AWS RDS CA bundle and split the bundle into individual certs (prefixed with cert)
# See http://blog.swwomm.com/2015/02/importing-new-rds-ca-certificate-into.html
ADD https://s3.amazonaws.com/rds-downloads/rds-combined-ca-bundle.pem /tmp/rds-ca/aws-rds-ca-bundle.pem
RUN cd /tmp/rds-ca && awk '/-BEGIN CERTIFICATE-/{close(x); x=++i;}{print > "cert"x;}' ./aws-rds-ca-bundle.pem \
    && for CERT in /tmp/rds-ca/cert*; do mv $CERT /usr/local/share/ca-certificates/aws-rds-ca-$(basename $CERT).crt; done \
    && rm -rf /tmp/rds-ca \
    && update-ca-certificates

COPY --from=build /go/src/github.com/dmlittle/discoverrewind/build .
COPY --from=build /go/src/github.com/dmlittle/discoverrewind/assets ./assets
COPY --from=build /go/src/github.com/dmlittle/discoverrewind/views ./views

RUN apk add --no-cache tzdata

CMD ["./serve"]
