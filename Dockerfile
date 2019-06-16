FROM golang:alpine as build
WORKDIR $GOPATH/src/github.com/TylerBrock/saw

# Add ca-certificates for TLS/SSL
RUN apk add --no-cache git ca-certificates

# Copy the rest of the project and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o /saw .

# Reset to scratch to drop all of the above layers and only copy over the final binary
FROM scratch
ENV HOME=/home
COPY --from=build /saw /bin/saw
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/bin/saw"]
