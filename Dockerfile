FROM golang:1.25-alpine AS builder

# install dependencies
RUN apk add --no-cache curl

# set working directory
WORKDIR /app

# copy the source code
COPY go.mod go.sum ./
RUN go mod download

# copy the rest of the source code
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/post-service cmd/post-service/main.go

FROM alpine:3.22

# install dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata

COPY --from=builder /app/post-service /app/post-service
RUN chmod +x /app/post-service

WORKDIR /app

ENV TZ=Asia/Jakarta

CMD ["/app/post-service"]


