FROM golang:1.13-alpine AS build-stage

RUN apk add --update make
RUN apk add --no-cache git

WORKDIR /httplie

COPY go.mod .
COPY go.sum .
COPY main.go .
RUN go get ./...

ARG GIT_COMMIT
ARG VERSION
ARG BUILD_DATE

COPY Makefile .
RUN make build

FROM alpine

RUN apk --update add ca-certificates
RUN mkdir /httplie
WORKDIR /httplie

COPY --from=build-stage  /httplie .

EXPOSE 8081

ENTRYPOINT ["./httplie"]
