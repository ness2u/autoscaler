FROM golang:latest as builder
WORKDIR /app
COPY . .
RUN go mod init autoscaler && \
    go get k8s.io/client-go@latest k8s.io/apimachinery@latest && \
    go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o autoscaler .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/autoscaler .
CMD ["./autoscaler"]