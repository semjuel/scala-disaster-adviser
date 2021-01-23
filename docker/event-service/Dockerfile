FROM golang:1.15-alpine

ENV SERVICE_PATH="scala-disaster-adviser/event-service"

RUN apk update && apk upgrade \
	&& apk add --no-cache ca-certificates \
	bash \
	git

RUN mkdir -p $GOPATH/src/${SERVICE_PATH}

# COPY source files
COPY ./event-service $GOPATH/src/${SERVICE_PATH}
COPY ./event-service/.env.example $GOPATH/src/${SERVICE_PATH}/.env

WORKDIR $GOPATH/src/${SERVICE_PATH}

# Install all dependencies.
RUN go mod vendor

CMD ["go", "run", "main.go"]
