FROM golang:alpine

WORKDIR temperature

COPY ./ ./

RUN ["go", "mod", "tidy"]
RUN ["go", "build", "./cmd/temperature/"]

ENTRYPOINT ["./temperature"]