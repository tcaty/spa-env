# -- build stage --
FROM golang:1.21.5-alpine3.18 as build

WORKDIR /usr/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./spa-env ./

# -- runtime stage --
FROM scratch as runtime

COPY --from=build /usr/src/spa-env /spa-env

ENTRYPOINT ["/spa-env"]
