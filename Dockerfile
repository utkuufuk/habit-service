FROM golang:alpine
WORKDIR /src
COPY go.sum go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /bin/app ./cmd/habit-service

CMD ["/bin/app"]
