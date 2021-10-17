FROM golang:alpine as build
RUN apk --no-cache add tzdata
WORKDIR /src
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/habit-service ./cmd/habit-service
RUN CGO_ENABLED=0 go build -o /bin/progress-report ./cmd/progress-report
RUN CGO_ENABLED=0 go build -o /bin/score-update ./cmd/score-update

FROM scratch
COPY --from=build /bin/habit-service /bin/habit-service
COPY --from=build /bin/progress-report /bin/progress-report
COPY --from=build /bin/score-update /bin/score-update
COPY --from=build /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /bin
CMD ["./habit-service"]
