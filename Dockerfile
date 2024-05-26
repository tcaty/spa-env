# -- Build stage --
FROM golang:1.21.5-alpine3.18 as build

WORKDIR /usr/src

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o ./spa-env ./

# -- Runtime stage --
FROM alpine:3.18 as runtime

# Install git in order to comfortable usage
# in ci systems runners. For example, gitlab runner.
RUN apk update && apk add --no-cache git 

COPY --from=build /usr/src/spa-env /spa-env

ENTRYPOINT ["/spa-env"]
