FROM golang:1.10.3-alpine3.8

RUN apk add git

COPY . /go/src/github.com/mahowald/telephone
WORKDIR /go/src/github.com/mahowald/telephone/cmd/telephone

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix nocgo

FROM python:3.7-alpine
WORKDIR /app
COPY --from=0 /go/src/github.com/mahowald/telephone/cmd/telephone/telephone /app/telephone

COPY test.py /app/

ENTRYPOINT ["./telephone"]