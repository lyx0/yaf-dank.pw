FROM golang:alpine as build
WORKDIR /app
COPY . .
RUN go build

FROM alpine:latest
COPY --from=build /app/yaf /app/yaf
WORKDIR /app
RUN mkdir -p /var/www/yaf
CMD ["./yaf"]
