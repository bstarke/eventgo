ARG goVersion=1.12
#ARG scannerVersion=0.0.14

# build stage
FROM hub.docker.prod.walmart.com/library/golang:${goVersion}-alpine3.10 AS build-env

WORKDIR /go/src/eventgo

COPY . .

RUN apk update && \
    apk add git && \
    export GIT_COMMIT=$(git rev-list -1 HEAD) && \
    export BUILD_TIME=$(date --utc +%FT%TZ) && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-X main.GitHash=$GIT_COMMIT -X main.GoVer=$GOLANG_VERSION -X main.BuildTime=$BUILD_TIME -w -s" -o app

#Non Root User Configuration
RUN addgroup -S -g 10001 appGrp \
    && adduser -S -D -u 10000 -s /sbin/nologin -h /app -G appGrp app \
    && chown -R 10000:10001 /app

# final stage
FROM scratch
LABEL maintainer="Brad Starkenberg"
# Import the user and group files from the builder.
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /go/src/eventgo/app .

#Override as non-root user
USER 10000:10001

EXPOSE 8080

ENTRYPOINT ["/app"]