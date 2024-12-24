FROM golang:1.23-alpine AS build
WORKDIR /go/src/app
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/sharex-server

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /go/bin/sharex-server .
EXPOSE 3939
ENTRYPOINT [ "/app/sharex-server" ]
