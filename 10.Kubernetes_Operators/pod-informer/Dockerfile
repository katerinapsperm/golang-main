##
## Build
##
FROM golang:1.17-buster AS build

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 go build -o /controller ./main.go

##
## Deploy
##
FROM scratch

WORKDIR /

COPY --from=build /controller /controller

ENTRYPOINT ["/controller"]