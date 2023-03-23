FROM golang:1.19.1 as build

LABEL author="vlbeaudoin"

WORKDIR /go/src/app

COPY go.mod go.sum main.go ./

ADD cmd/ cmd/

ADD api/ api/

RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o haul .

# Alpine

FROM alpine:3.17.2

RUN apk update && apk add file --no-cache && apk upgrade --no-cache

WORKDIR /app

COPY --from=build /go/src/app/haul /usr/bin/haul

CMD ["haul", "--help"]
