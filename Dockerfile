FROM golang:alpine as build

LABEL maintainer="Leon Jacobs <leonja511@gmail.com>"

COPY . /src

WORKDIR /src
RUN go build -o nutstat

# final image
FROM golang:alpine

COPY --from=build /src/nutstat /

VOLUME ["/"]

ENTRYPOINT ["/nutstat"]
