FROM golang:alpine

WORKDIR /go/src/github.com/go-park-mail-ru/2018_2_parashutnaya_molitva

COPY . .

RUN go install github.com/go-park-mail-ru/2018_2_parashutnaya_molitva/cmd/server

CMD ["server"]