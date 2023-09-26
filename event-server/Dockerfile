FROM golang:1.19-alpine as builder

# Setup
RUN mkdir -p /app
WORKDIR /app

# Add libraries
RUN apk add --no-cache git

# Copy & build
ADD ./event-server/* /app
RUN mkdir build 
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -o build/event-server ./main.go

# Copy into scratch container
FROM alpine
COPY --from=builder /app/build/event-server ./
ENTRYPOINT ["./event-server"]