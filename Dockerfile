# Build the image
FROM golang:alpine AS build
RUN apk add build-base
COPY . /src
WORKDIR /src
RUN go build

# Build the final image
FROM alpine:edge
COPY --from=build /src/shoutbox-server /shoutbox-server

# Start the server
WORKDIR /
ENTRYPOINT [ "/shoutbox-server" ]
EXPOSE 8000
