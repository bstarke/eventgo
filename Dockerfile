# build stage
FROM harbor.home.starkenberg.net/hub/golang:1.22 AS build-env

WORKDIR /go/src/eventgo

COPY . .

RUN apt update && \
    apt install git && \
    export GIT_COMMIT=$(git rev-list -1 HEAD) && \
    export BUILD_TIME=$(date --utc +%FT%TZ) && \
    GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-X main.GitHash=$GIT_COMMIT -X main.GoVer=$GOLANG_VERSION -X main.BuildTime=$BUILD_TIME -w -s" -o app

#Non Root User Configuration
RUN adduser --system --group --disabled-login --shell /sbin/nologin --home /app/ app

# final stage
FROM scratch
LABEL maintainer="Brad Starkenberg"
# Import the user and group files from the builder.
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /go/src/eventgo/app /app/app

#Override as non-root user
USER app

EXPOSE 8080

ENTRYPOINT ["/app/app"]