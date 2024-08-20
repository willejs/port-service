# Build stage
FROM golang:1.22 AS builder

# do people even care about the LFS hirarchy anymore?
WORKDIR /app

# Copy go mod and sum files in first so we can cache the dependencies
COPY go.mod go.sum ./
RUN go mod download

# copy in the app and build it
COPY . .

# statically compile the go binary for the presumed target of amd64 linux
# whilst its larger, its more portable and will run in a scratch container
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o api-server ./cmd/api-server

# hack to create the nobody user for the scratch container.
# hadolint ignore=DL3059
RUN echo "nobody:x:65534:65534:Nobody:/:" > /etc_passwd

# Use multi stage builds. This is the final runtime stage. We can use ephemeral containers in kubernetes now :tada:
# We could make a development target if this displeases people too.
FROM scratch

# expose the port and hardcode it for now
EXPOSE 8080

# copy in the nobody user in
COPY --from=builder /etc_passwd /etc/passwd

WORKDIR /app
COPY --from=builder /app/api-server /app/api-server

# dont run the app as root, it is insecure
USER nobody
CMD ["./api-server"]