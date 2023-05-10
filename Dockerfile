# syntax=docker/dockerfile:1

FROM golang:1.20 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /proxy -a -ldflags '-linkmode external -extldflags "-static"' .


FROM scratch
COPY --from=builder /proxy /proxy
COPY *.db ./

EXPOSE 8080

CMD ["/proxy", "-run"]
