FROM golang:1.19-alpine3.16 AS build
WORKDIR /app
COPY . .

RUN  CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app .

FROM alpine:latest
COPY --from=build /app/app /app
COPY private-key.pem .
CMD [ "/app" ]
