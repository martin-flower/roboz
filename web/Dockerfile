# docker build --tag roboz-web:minimum --file Dockerfile .
# docker run --detach --publish 5000:5000 roboz-web:minimum
# docker run --publish 5000:5000 roboz-web:minimum

# multistage build based on example in https://raw.githubusercontent.com/cloud-native-go/examples/main/ch05/ch05_0X_dockerbuild/Dockerfile

# 1 -- -- -- go tests
FROM golang:1.20.4 as test

WORKDIR /testdir

# Fetch dependencies
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# RUN go build -o main .

# RUN go mod tidy

# run the go unittests - a failing test will fail the build
RUN go test -timeout 90s -v ./...


# 2 -- -- -- compile go source to binary
FROM golang:1.20.4 as build

# copy the source files from the host
COPY . /src

# copy the dependency source from the test container so we don't have to re-download it
COPY --from=test /go/pkg/mod/ /go/pkg/mod/

# Set the working directory to the same place we copied the code
WORKDIR /src

# build the binary
#  CGO_ENABLED=0 - do not use CGO; compile statically
#  GOOS=linux    - compile for linux os
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o roboz-web

# 3 -- -- -- build image

# scratch image contains no distribution files
# resulting image and containers will have only the roboz-web binary
FROM scratch as image

# copy the binary from the build container
COPY --from=build /src/roboz-web .

# copy the configuration from the build container
COPY --from=build /src/config/roboz.yaml .

# http port
EXPOSE 5000

# execute this with docker run
CMD ["/roboz-web"]
