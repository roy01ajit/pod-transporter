# Build the manager binary
FROM golang:1.18 as builder

WORKDIR /pod-transporter
# Copy the Go Modules manifests
COPY . .
# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o main ./apis/cmd


FROM alpine:latest As production
WORKDIR /root/
# Copy the Pre-built binary file from the previous stage
COPY --from=builder /pod-transporter/main .

EXPOSE 9090

# Command to run the executable
CMD ["./main"]
