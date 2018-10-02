FROM golang:alpine

COPY ./../../. /go/src/github.com/go-park-mail-ru/2018_2_parashutnaya_molitva

RUN  install github.com/go-park-mail-ru/2018_2_parashutnaya_molitva

CMD ["go", "run", "/go/src/github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/server/server.go"]