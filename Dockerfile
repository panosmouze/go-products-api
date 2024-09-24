FROM golang:1.23-alpine AS build
RUN apk add --no-cache gcc musl-dev
ENV CGO_ENABLED=1
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -o myapp .

FROM alpine:latest
ENV APP_ENV=production
RUN apk add --no-cache libc6-compat
WORKDIR /app
COPY --from=build /app/myapp .
EXPOSE 8080
CMD ["./myapp"]