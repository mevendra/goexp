FROM golang:alpine

WORKDIR limiter

COPY ./ ./

RUN ["go", "mod", "tidy"]
RUN ["go", "build", "./cmd/limiter/"]

ENTRYPOINT ["./limiter"]