FROM golang:alpine AS build
RUN apk --no-cache add gcc g++ make git
WORKDIR /go/src/app
COPY . .
# COPY ./ducksdockside.com /go/src/bytesupply.com
RUN go get ./...
RUN GOOS=linux go build -ldflags="-s -w" -o ./bin/duckysdockside-app ./main.go
RUN mkdir ./bin/log
# RUN mkdir -p ./bin/data/qTurHm/
COPY ./static/ ./bin/static/
COPY ./static/html/ ./bin/static/html/
COPY ./sitemap.xml ./bin/sitemap.xml
COPY ./robots.txt ./bin/robots.txt
RUN ["chmod", "+x", "/bin"]
FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /usr/bin
COPY --from=build /go/src/app/bin /go/bin
EXPOSE 80
ENTRYPOINT /go/bin/duckysdockside-app --port 80
