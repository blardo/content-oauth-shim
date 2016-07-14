# To build:
# $ docker run --rm -v $(pwd):/go/src/github.com/skuid/content-oauth-shim -w /go/src/github.com/skuid/content-oauth-shim golang:1.6  go build -v -a -tags netgo -installsuffix netgo -ldflags '-w'
# $ docker build -t skuid/content-oauth-shim .
#
# To run:
# $ docker run skuid/content-oauth-shim

FROM busybox

MAINTAINER Micah Hausler, <micah@skuid.com>

COPY content-oauth-shim /bin/content-oauth-shim
RUN chmod 755 /bin/content-oauth-shim

EXPOSE 3000

ENTRYPOINT ["/bin/content-oauth-shim"]
