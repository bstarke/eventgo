ARG goVersion=1.12
#ARG scannerVersion=0.0.14

# build stage
#FROM hub.docker.prod.walmart.com/library/golang:${goVersion}-alpine3.8 AS build-env
FROM golang:${goVersion}-alpine3.10 AS build-env
WORKDIR /go/src/eventgo
COPY . .
RUN ls -lah
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags="-w -s" -o app

#Non Root User Configuration
RUN addgroup -S -g 10001 appGrp \
    && adduser -S -D -u 10000 -s /sbin/nologin -h /app -G appGrp app \
    && chown -R 10000:10001 /app

# test stage
#FROM build-env as test
#RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go test ./... -coverprofile cover.out -json > report.json

#sonar
#FROM docker.prod.walmart.com/strati/docker-sonarqube-scanner:${scannerVersion} as sonar
#ARG goSourceDir
#ARG sonarOpts
#ARG sonarProjKey
#COPY --from=test ${goSourceDir} /app
#WORKDIR /app
#RUN sonar-scanner --debug -Dsonar.projectKey=${sonarProjKey} ${sonarOpts}

# final stage
FROM scratch
LABEL maintainer="Walmart Container Engineering"
# Import the user and group files from the builder.
COPY --from=build-env /etc/passwd /etc/passwd
COPY --from=build-env /go/src/eventgo/app .

#Override as non-root user
USER 10000:10001

EXPOSE 8080

ENTRYPOINT ["/app"]