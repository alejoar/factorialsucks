ARG APP=factorialsucks

# Define Go version
ARG GO_VERSION=1.18

# Build container
FROM docker.io/golang:${GO_VERSION}-alpine AS build
ARG APP

RUN apk add --no-cache git

WORKDIR /src

# Fetch dependencies
COPY ./go.mod ./go.sum ./
RUN go mod download

COPY ./ ./

# Build static binary
RUN CGO_ENABLED=0 go build -installsuffix 'static' -o /${APP} .

# Final container
FROM gcr.io/distroless/static AS final
ARG APP

USER nonroot:nonroot

COPY --from=build --chown=nonroot:nonroot /${APP} /executable

ENTRYPOINT ["/executable"]
