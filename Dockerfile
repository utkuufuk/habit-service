FROM golang:alpine as build
RUN apk --no-cache add tzdata
WORKDIR /src
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/habit-service

FROM scratch
COPY --from=build /bin/app /bin/app
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /bin
CMD ["./app"]
