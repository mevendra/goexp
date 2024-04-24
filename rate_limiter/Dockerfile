FROM golang:alpine

WORKDIR limiter

COPY ./ ./
COPY ./cmd/limiter/env ./

RUN ["go", "mod", "tidy"]
RUN ["go", "build", "./cmd/limiter/"]

ENTRYPOINT ["./limiter"]