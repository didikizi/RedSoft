FROM golang:alpine3.20 as builder

LABEL version="1.0.0"

ENV PATH_PROJECT=/app
ENV GO111MODULE=on
ENV GOSUMDB=off
ENV TARGET=testTask

WORKDIR ${PATH_PROJECT}
COPY . ${PATH_PROJECT}
RUN CGO_ENABLED=0 GOOS=linux go build ./cmd/$TARGET/

FROM alpine:3.20
ENV PATH_PROJECT=/app
ENV TARGET=testTask
COPY --from=builder $PATH_PROJECT/$TARGET /bin
COPY --from=builder $PATH_PROJECT/iternal/migrations iternal/migrations
CMD ["/bin/testTask"]
EXPOSE  4000