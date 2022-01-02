# docker build --tag roboz-smoke:latest .
# docker run -it --network=host roboz-smoke:latest

FROM golang:1.17.5

WORKDIR /smoke

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENTRYPOINT ["go", "test", "--verbose"]
