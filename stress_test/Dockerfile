FROM golang:alpine

WORKDIR stress_test
COPY ./ ./
RUN ["go", "mod", "tidy"]
RUN ["go", "build", "."]

ENTRYPOINT ["./stress_test"]