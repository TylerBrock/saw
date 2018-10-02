FROM golang:alpine as build
WORKDIR $GOPATH/src/github.com/TylerBrock/saw

# Setup some basic dependencies that aren’t bundled in the build image
ADD https://github.com/golang/dep/releases/download/v0.5.0/dep-linux-amd64 /bin/dep
RUN chmod +x /bin/dep
RUN apk add --no-cache git

# Ensure deps separately for a cache layer during rebuilds
COPY Gopkg.* ./
RUN dep ensure -vendor-only

# Copy the rest of the project and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /saw .

# Reset to scratch to drop all of the above layers and only copy over the final binary
FROM scratch
COPY --from=build /saw /bin/saw
ENTRYPOINT ["/bin/saw"]
