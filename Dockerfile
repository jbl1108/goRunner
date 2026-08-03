# base image
FROM golang:1.25-alpine

# define the working directory
WORKDIR /app

# copy the go.mod and go.sum so that the packages to be installed
# are known in the container. ./ here is the WORKDIR, /app
COPY go.mod ./

# command to install modules
RUN go mod download

# copy source code into working dir
COPY ./. .
#expose the port the app runs on
EXPOSE 9090

# build
RUN CGO_ENABLED=0 GOOS=linux go build -o /goRunner ./goRunner.go

VOLUME ["/gorunner/config"]

# run the compiled binary when the container starts
CMD ["/goRunner"]