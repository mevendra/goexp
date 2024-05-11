FROM golang:alpine

WORKDIR ordersystem

COPY ./ ./

RUN ["go", "mod", "tidy"]
RUN ["go", "build", "./cmd/ordersystem/"]

ENTRYPOINT ["./ordersystem"]