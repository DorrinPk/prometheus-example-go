FROM golang AS build
WORKDIR /src

COPY metrics.go ./
RUN go build metrics.go

FROM golang
WORKDIR /src

COPY --from=build /src/metrics /usr/local/bin/

COPY server.go ./
COPY go.mod ./
COPY go.sum ./
RUN go build server.go

EXPOSE 8080
CMD ["./server"]
