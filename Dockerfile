############################
# STEP 1: Build the binary.
############################
FROM golang:1.14.7-alpine AS builder

# Install git.
# Git is required for fetching the dependencies.
RUN apk update && apk add --no-cache git openssh-client

# Import username and password.
ARG ghuser
ARG ghpass

# Ensure SSH is used to access GitHub.
RUN git config --global url."https://${ghuser}:${ghpass}@github.com/".insteadOf "https://github.com/"

# Copy the code to build.
WORKDIR /api
COPY . .

# Build the binary.
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-w -s" -o /go/bin/api

############################
# STEP 2: Transfer to minimal image.
############################
FROM alpine

# Copy trusted certificates.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Copy our static executable.
COPY --from=builder /go/bin/api /go/bin/api

# Fusl recommended these and I trust her.
ENV GODEBUG=madvdontneed=1
ENV GOGC=20

# Run the hello binary.
ENTRYPOINT ["/go/bin/api"]