FROM golang:alpine as build
WORKDIR /go/src/github.com/Ergotu/Amass
COPY . .
RUN apk --no-cache add git \
  && go get -u github.com/Ergotu/Amass/...
  
FROM alpine:latest
COPY --from=build /go/bin/amass /bin/amass 
ENTRYPOINT ["/bin/amass"]
