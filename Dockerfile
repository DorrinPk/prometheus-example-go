FROM golang AS build
WORKDIR /src

COPY metrics.go ./
RUN go build metrics.go

FROM golang
WORKDIR /src

COPY --from=build /src/metrics /usr/local/bin/

# Build your code here, the provided server is an example and can be
# replaced with any other desired language. You can change the base
# image as well if it helps in installation or testing.
COPY server.go ./
COPY go.mod ./
COPY go.sum ./
RUN go build server.go

EXPOSE 8080
CMD ["./server"]
