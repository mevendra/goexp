FROM golang:alpine

WORKDIR temperature

COPY ./ ./

RUN ["go", "mod", "tidy"]
RUN ["go", "build", "./cmd/temperature/"]

EXPOSE 8080

ENTRYPOINT ["./temperature"]