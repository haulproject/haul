FROM golang:1.20.2 as build

LABEL author="vlbeaudoin"

WORKDIR /go/src/app

COPY go.mod go.sum main.go ./

ADD cmd/ cmd/

ADD api/ api/

ADD cli/ cli/

ADD types/ types/

ADD db/ db/

ADD handlers/ handlers/

ADD graph/ graph/

#RUN CGO_ENABLED=0 go build -a -installsuffix cgo -o haul .
RUN go build -a -installsuffix cgo -o haul .

# Debian

FROM debian:stable-20230502

WORKDIR /app

COPY --from=build /go/src/app/haul /usr/bin/haul

CMD ["haul", "--help"]
