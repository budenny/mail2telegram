FROM golang:1.15-alpine AS builder
RUN apk add --no-cache git
WORKDIR /build

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

RUN go build -o /build/mail2telegram .

###########################################################

FROM alpine:3.13.2 
COPY --from=builder /build/mail2telegram /app/mail2telegram
CMD ["/app/mail2telegram"]