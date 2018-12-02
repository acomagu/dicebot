FROM golang:alpine as build-env
RUN apk --update add ca-certificates git
ADD . /src
ENV CGO_ENABLED=0
RUN cd /src && go build -o dicebot

FROM scratch
WORKDIR /app
ENV PATH=/bin
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build-env /src/dicebot /app/
ENTRYPOINT ["./dicebot"]
